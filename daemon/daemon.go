package main

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

// RPCDataFetcher defines a subset of methods conformed to by ETH1.0 RPC clients for
// fetching eth1 data from the clients.
type RPCDataFetcher interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*gethTypes.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*gethTypes.Header, error)
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
}

type RPCClient interface {
	BatchCall(b []gethRPC.BatchElem) error
}

type Service struct {
	connected       bool
	isRunning       bool
	cfg             *Web3ServiceConfig
	ctx             context.Context
	cancel          context.CancelFunc
	httpLogger      bind.ContractFilterer
	eth1DataFetcher RPCDataFetcher
	rpcClient       RPCClient
	runError                error
}

// Web3ServiceConfig defines a config struct for web3 service to use through its life cycle.
type Web3ServiceConfig struct {
	HttpEndpoints      string
	Eth1HeaderReqLimit uint64
}

func NewService(ctx context.Context) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop()
	s := &Service{
		ctx:    ctx,
		cancel: cancel,
	}
	return s, nil
}

// Start a web3 service's main event loop.
func (s *Service) Start() {
	go func() {
		s.isRunning = true
		s.waitForConnection()
		if s.ctx.Err() != nil {
			log.Info("Context closed, exiting pow goroutine")
			return
		}
		s.run(s.ctx.Done())
	}()
}

func (s *Service) Stop() error {
	if s.cancel != nil {
		defer s.cancel()
	}
	s.closeClients()
	return nil
}


// run subscribes to all the services for the ETH1.0 chain.
func (s *Service) run(done <-chan struct{}) {
	s.runError = nil

	s.initPOWService()
}

func (s *Service) initPOWService() {

	// Run in a select loop to retry in the event of any failures.
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			ctx := s.ctx
			header, err := s.eth1DataFetcher.HeaderByNumber(ctx, nil)
			if err != nil {
				log.Errorf("Unable to retrieve latest ETH1.0 chain header: %v", err)
				continue
			}

			log.Infof("Starting pow service with ETH1.0 header: %v", header)
			return
		}
	}
}

func (s *Service) waitForConnection() {
	errConnect := s.connectToPowChain()
	if errConnect == nil {
		log.Info("Connected to eth1 nodes")
	}

	if errConnect != nil {
		s.runError = errConnect
		log.Errorf("Could not connect to powchain endpoint")
	}
}

func (s *Service) connectToPowChain() error {
	httpClient, rpcClient, err := s.dialETH1Nodes("ws://localhost:8576")
	if err != nil {
		return errors.Wrap(err, "could not dial eth1 nodes")
	}

	if httpClient == nil || rpcClient == nil {
		return errors.New("eth1 client is nil")
	}

	s.initializeConnection(httpClient, rpcClient)
	return nil
}

func (s *Service) dialETH1Nodes(url string) (*ethclient.Client, *gethRPC.Client, error) {
	httpRPCClient, err := gethRPC.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	httpClient := ethclient.NewClient(httpRPCClient)
	// Add a method to clean-up and close clients in the event
	// of any connection failure.
	closeClients := func() {
		httpRPCClient.Close()
		httpClient.Close()
	}

	syncProg, err := httpClient.SyncProgress(s.ctx)
	if err != nil {
		closeClients()
		return nil, nil, err
	}
	
	if syncProg != nil {
		closeClients()
		return nil, nil, errors.New("eth1 node has not finished syncing yet")
	}
	
	// Make a simple call to ensure we are actually connected to a working node.
	cID, err := httpClient.ChainID(s.ctx)
	if err != nil {
		closeClients()
		return nil, nil, err
	}
	nID, err := httpClient.NetworkID(s.ctx)
	if err != nil {
		closeClients()
		return nil, nil, err
	}
	log.Infof("ChainID: %s, NetworkID: %s", cID.String(), nID.String())

	headers := make(chan *types.Header)
	sub, err := httpClient.SubscribeNewHead(s.ctx, headers)
	if err != nil {
		log.Errorln("Cannot SubscribeNewHead to host. Error: ", err)
		return nil, nil, err
	}

	go func() {
		for {
			select {
						case err := <-sub.Err():
							log.Errorln("Error: ", err)
						case header := <-headers:
							log.Infoln("Block Number: ", header.Number.String())
				}
		}
	}()

	return httpClient, httpRPCClient, nil
}

// closes down our active eth1 clients.
func (s *Service) closeClients() {
	gethClient, ok := s.rpcClient.(*gethRPC.Client)
	if ok {
		gethClient.Close()
	}
	httpClient, ok := s.eth1DataFetcher.(*ethclient.Client)
	if ok {
		httpClient.Close()
	}
}

func (s *Service) initializeConnection(
	httpClient *ethclient.Client,
	rpcClient *gethRPC.Client,
) {
	s.httpLogger = httpClient
	s.eth1DataFetcher = httpClient
	s.rpcClient = rpcClient
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("Hello, World!")
	service, err := NewService(context.Background())
	if err != nil {
		log.Errorln("Cannot initialize Service")
	}
	service.Start()
	wg.Wait()
}