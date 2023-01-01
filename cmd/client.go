package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anvari1313/grpc-load/internal/client"
)

var (
	clientCMD = &cobra.Command{
		Use:   "client",
		Short: "A generator for Cobra based Applications",
		Run:   clientFunc,
	}
)

func clientFunc(cmd *cobra.Command, args []string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	serverAddress, found := os.LookupEnv("GRPC_LOAD_SERVER_ADDRESS")
	if !found {
		serverAddress = "127.0.0.1:1313"
	}

	client.Init(":1213", serverAddress, logger)
}
