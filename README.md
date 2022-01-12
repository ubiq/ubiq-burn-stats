# Watch the Burn 🔥
When EIP-1559 gets deployed, ETH will be burned in every block if transactions exist. This website will show you how much ETH got burned in total and per block.

If you have a local Ubiq (we use gubiq but you can use any ETH client) instance, you can update `REACT_APP_WEB3_URL` in  `.env.production.local` to your local ETH instance and run this offline! The instructions below show how to deploy it to a remote website under nginx.

[![Frontend CI/CD](https://github.com/ubiq/ubiq-burn-stats/actions/workflows/frontend-azure-static-web-apps.yml/badge.svg?branch=main)](https://github.com/ubiq/ubiq-burn-stats/actions/workflows/frontend-azure-static-web-apps.yml) [![Daemon CI/CD](https://github.com/ubiq/ubiq-burn-stats/actions/workflows/daemon-linode.yml/badge.svg?branch=main)](https://github.com/ubiq/ubiq-burn-stats/actions/workflows/daemon-linode.yml)

## Setup dev environment ⚙

Setting up the environment requires a gubiq instance, daemon gubiq proxy, and react web app. Optionally, you can install the varnish cache in between.

### Setup gubiq
1. Clone gubiq and build docker image. Assumes `/data` on local system exists
   ```
   git clone https://github.com/ubiq/go-ubiq.git
   cd go-ubiq
   docker build -t ubiq-node .
   mkdir /data
   ```

1. To run Gubiq inside Docker, run one of the following:
   *  Mainnet
      ```
       docker run --net=host --name=gubiq -dti -v /data:/data ubiq-node --datadir=/data/mainnet --mainnet --port=3030 --http --http.port=8588 --http.api="net,web3,eth" --http.corsdomain="localhost"  --ws --ws.port=8589 --ws.api="net,web3,eth" --ws.origins="*" --maxpeers=100
      ```
### Setup the Gubiq Proxy Daemon

1. Easiest thing is use docker.
   ```
    docker build -t gubiq-proxy ./daemon
   ```
   
1. Run the docker instance against your gubiq.
   ```
   docker run -d --name=gubiq-proxy --restart=on-failure:3 --net=host -v /data/gubiq-proxy:/data gubiq-proxy --addr=:8080 --gubiq-endpoint-http=http://localhost:8588 --gubiq-endpoint-websocket=ws://localhost:8589 --db-path=/data/mainnet.db --development
   ```
   If you include `--initializedb` it will start initializing the database since EIP-London, will take time. If you take it out, then it basically just starts at the current head.
   
### Optional: Varnish cache to cache all Gubiq RPC calls

1. Easiest thing is use docker.
   ```
    docker build -t gubiq-varnish ./cache
   ```
   
1. Run the docker instance against your gubiq.
   ```
   docker run -d --name gubiq-cache --net host -e GETH_HTTP_HOST=localhost -e GETH_HTTP_PORT=8588 -e VARNISH_PORT=8081 -e CACHE_SIZE=1g gubiq-varnish
   ```

1. Run the daemon against the varnish port.
   ```
   docker run -d --name=gubiq-proxy --restart=on-failure:3 --net=host -v /data/gubiq-proxy:/data gubiq-proxy --addr=:8080 --gubiq-endpoint-http=http://localhost:8081 --gubiq-endpoint-websocket=ws://localhost:8589 --db-path=/data/mainnet.db --development
   ```
   
### Setup web dev environment

1. Create `env` file:
   ```
   cp .env .env.local
   ```

1. Add your gubiq ws url to `.env.local` point it to your daemon port:
   ```
   REACT_APP_WEB3_URL=localhost:8080
   ```

1. Install packages
   ```
   npm install
   ```

1. Run the web app:
   ```
   npm start
   ```
1. Launch the web app (goerli):
   ```
   open http://mainnet.go.localhost:3000:
   ```

### Some devops maintenance

**Send transactions to testnet** 

Install web3 CLI client `curl -LSs https://raw.githubusercontent.com/gochain/web3/master/install.sh | sh` use it to create a test account, and you can use it to send transactions.

**Access gubiq console**
If you ran the `mainnet` docker gubiq, you can just do, not `rm` is there so it cleans the container up after closing!:
```
docker run --rm -ti -v /data:/data ubiq-node --datadir=/data/mainnet attach  
```
