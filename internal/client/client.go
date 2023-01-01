package client

import (
	"context"
	"net/http"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/anvari1313/grpc-load/proto"
)

func Init(bind, gRPCAddress string, logger *zap.Logger) {
	logger.Info("dialing gRPC client", zap.String("server_address", gRPCAddress))

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	}
	conn, err := grpc.Dial(gRPCAddress, opts...)
	if err != nil {
		logger.Fatal("error in dialing gRPC client", zap.Error(err))
	}
	defer conn.Close()

	client := proto.NewGreeterClient(conn)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status))
			return nil
		},
	}))

	e.GET("/", runAllMethodsSyncHandler(client))
	e.GET("/batch", runBatchAsyncHandler(client, logger))

	logger.Info("starting HTTP server", zap.String("bind", bind))
	e.Logger.Fatal(e.Start(bind))
}

func runAllMethodsSyncHandler(gRPCClient proto.GreeterClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := gRPCClient.SayHello(c.Request().Context(), &proto.HelloRequest{Name: "this"})
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, result.Message)
	}
}

type batchResponse struct {
	Count uint64        `json:"count"`
	Delay time.Duration `json:"delay"`
}

func runBatchAsyncHandler(gRPCClient proto.GreeterClient, logger *zap.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		count := uint64(1500)
		delay := 1 * time.Millisecond

		go func() {
			for i := uint64(0); i < count; i++ {
				_, err := gRPCClient.SayHello(context.Background(), &proto.HelloRequest{Name: "this"})
				if err != nil {
					logger.Error("error in calling SayHello", zap.Error(err))
				}

				time.Sleep(delay)
			}
		}()

		return c.JSON(http.StatusOK, batchResponse{
			Count: count,
			Delay: delay,
		})
	}
}
