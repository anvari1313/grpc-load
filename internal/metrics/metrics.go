package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func Serve(bind string, logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Info("starting metrics HTTP server", zap.String("bind", bind))
		err := http.ListenAndServe(bind, nil)
		if err != nil {
			logger.Fatal("error in initiating metrics HTTP server", zap.Error(err))
		}
	}()
}
