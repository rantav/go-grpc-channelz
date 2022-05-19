package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	demoservice "github.com/rantav/go-grpc-channelz/internal/generated/service"
)

type server struct {
	demoservice.UnimplementedDemoServiceServer
}

// New creates a new grpc server
func New() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	server := &server{}
	demoservice.RegisterDemoServiceServer(grpcServer, server)
	return grpcServer, nil
}

func (s *server) Hello(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
