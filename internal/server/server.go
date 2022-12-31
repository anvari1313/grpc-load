package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/anvari1313/grpc-load/proto"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "this is message"}, nil
}

func (s *Server) mustEmbedUnimplementedGreeterServer() {
	//TODO implement me
	panic("implement me")
}

func StartServer() {
	lis, err := net.Listen("tcp", ":1313")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterGreeterServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
