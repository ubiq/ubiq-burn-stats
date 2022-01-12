package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ubiq/ubiq-burn-stats/daemon/hub"
)

func newRootCmd() *cobra.Command {
	var addr string
	var debug bool
	var development bool
	var gethEndpointHTTP string
	var gethEndpointWebsocket string
	var dbPath string
	var workerCount int

	rootCmd := &cobra.Command{
		// TODO:
		Use: "ubiq-burn-stats",
		// TODO:
		Short: "short",
		Long:  `long`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if gethEndpointHTTP == "" {
				cmd.Help()
				return fmt.Errorf("--gubiq-endpoint-http is required")
			}

			if gethEndpointWebsocket == "" {
				cmd.Help()
				return fmt.Errorf("--gubiq-endpoint-websocket is required")
			}

			if dbPath == "" {
				cmd.Help()
				return fmt.Errorf("--gubiq-endpoint-websocket is required")
			}

			return root(
				addr,
				debug,
				development,
				gethEndpointHTTP,
				gethEndpointWebsocket,
				dbPath,
				workerCount,
			)
		},
	}

	rootCmd.Flags().StringVar(&addr, "addr", ":8080", "HTTP service address")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "enable debug logs")
	rootCmd.Flags().BoolVar(&development, "development", true, "enable for development mode")
	rootCmd.Flags().StringVar(&gethEndpointHTTP, "gubiq-endpoint-http", "http://localhost:8588", "Endpoint to gubiq for http")
	rootCmd.Flags().StringVar(&gethEndpointWebsocket, "gubiq-endpoint-websocket", "ws://localhost:8589", "Endpoint to gubiq for websocket")
	rootCmd.Flags().StringVar(&dbPath, "db-path", "watchtheburn.db", "Path to the SQLite db")
	rootCmd.Flags().IntVar(&workerCount, "worker-count", 10, "Number of workers to spawn to parallelize http client")

	return rootCmd
}

func root(
	addr string,
	debug bool,
	development bool,
	gethEndpointHTTP string,
	gethEndpointWebsocket string,
	dbPath string,
	workerCount int,
) error {
	hub, err := hub.New(
		debug,
		development,
		gethEndpointHTTP,
		gethEndpointWebsocket,
		dbPath,
		workerCount,
	)
	if err != nil {
		return err
	}

	err = hub.ListenAndServe(addr)
	if err != nil {
		return err
	}

	return nil
}

// Execute runs cli command
func Execute() error {
	return newRootCmd().Execute()
}
