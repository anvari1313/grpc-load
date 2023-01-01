package server

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/anvari1313/grpc-load/proto"
)

type Server struct {
	proto.UnimplementedGreeterServer

	grpcServer *grpc.Server
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "this is message"}, nil
}

func InitUnary(bind string, logger *zap.Logger) *Server {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		logger.Fatal("error in listening on gRPC bind address", zap.String("bind", bind))
	}

	s := new(Server)

	opts := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_recovery.UnaryServerInterceptor(),
	))

	s.grpcServer = grpc.NewServer(opts)
	proto.RegisterGreeterServer(s.grpcServer, s)

	go func() {
		logger.Info("starting gRPC server", zap.String("bind", bind))
		err = s.grpcServer.Serve(lis)
		if err != nil {
			logger.Fatal("error in initiating gRPC server", zap.Error(err))
		}
	}()

	return s
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
