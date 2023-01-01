package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anvari1313/grpc-load/internal/metrics"
	"github.com/anvari1313/grpc-load/internal/server"
)

var (
	serverCMD = &cobra.Command{
		Use:   "server",
		Short: "A generator for Cobra based Applications",
		Run:   serverFunc,
	}
)

func serverFunc(cmd *cobra.Command, args []string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	metrics.Serve(":1212", logger)
	s := server.InitUnary(":1313", logger)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh

	s.Shutdown()
}
