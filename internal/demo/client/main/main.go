package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "google.golang.org/grpc/grpclog"

	"github.com/rantav/go-grpc-channelz/internal/demo/client"
)

func main() {
	const dialString = "localhost:8080"
	client, err := client.New(dialString)
	if err != nil {
		log.Fatalf("Cannot create gRPC client to %s. %v", dialString, err)
	}
	_, err = client.Hello(context.Background(), &empty.Empty{})
	if err != nil {
		log.Errorf("Error saying hello. %+v", err)
		return
	}
	time.Sleep(10 * time.Second)
	fmt.Println("Hello was successful")
}
