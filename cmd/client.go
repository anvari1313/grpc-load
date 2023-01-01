package cmd

import (
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

	client.Init(":1213", ":1313", logger)
}
