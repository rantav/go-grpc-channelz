package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	demoservice "github.com/rantav/go-grpc-channelz/internal/generated/service"
)

type server struct{}

// New creates a new grpc server
func New() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	server := &server{}
	demoservice.RegisterDemoServiceServer(grpcServer, server)
	return grpcServer, nil
}

func (s *server) Hello(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
