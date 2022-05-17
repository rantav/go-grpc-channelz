package client

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	demoservice "github.com/rantav/go-grpc-channelz/internal/generated/service"
)

// New creates a new gRPC client
func New(connectionString string) (demoservice.DemoServiceClient, error) {
	conn, err := grpc.Dial(connectionString,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error doaling to %s", connectionString)
	}

	client := demoservice.NewDemoServiceClient(conn)
	return client, err
}
