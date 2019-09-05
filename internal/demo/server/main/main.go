package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
	channelzservice "google.golang.org/grpc/channelz/service"
	log "google.golang.org/grpc/grpclog"

	"github.com/golang/protobuf/ptypes/empty"
	channelz "github.com/rantav/go-grpc-channelz"
	"github.com/rantav/go-grpc-channelz/internal/demo/client"
	"github.com/rantav/go-grpc-channelz/internal/demo/server"
)

func main() {
	const (
		grpcBindAddress  = ":8080"
		adminBindAddress = ":8081"
	)

	// nolint:gosec
	grpcListener, err := net.Listen("tcp", grpcBindAddress)
	if err != nil {
		log.Fatal(err)
	}

	// nolint:gosec
	adminListener, err := net.Listen("tcp", adminBindAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer, err := server.New()
	if err != nil {
		log.Fatalf("Failed to create grpc server %+v", err)
	}

	// Register the channelz handler
	http.Handle("/", channelz.CreateHandler("/foo", grpcBindAddress))

	// Register the channelz gRPC service to grpcServer so that we can query it for this service.
	channelzservice.RegisterChannelzServiceToServer(grpcServer)

	g := new(errgroup.Group)
	g.Go(func() error { return http.Serve(adminListener, nil) })
	g.Go(func() error { return grpcServer.Serve(grpcListener) })

	fmt.Printf("demo server is up is up; gRPC bind address: %s, http admin address: %s \n",
		grpcBindAddress, adminBindAddress)

	go runClient(fmt.Sprintf("localhost%s", grpcBindAddress))

	// should never return
	err = g.Wait()
	log.Fatalf("Error running server: %+v", err)
}

// runs a client gRPC call in a loop with some sleeps.
func runClient(dialString string) {
	client, err := client.New(dialString)
	if err != nil {
		log.Fatalf("Cannot create gRPC client to %s. %v", dialString, err)
	}
	for {
		time.Sleep(10 * time.Second)
		_, err = client.Hello(context.Background(), &empty.Empty{})
		if err != nil {
			log.Errorf("Error saying hello. %+v", err)
			return
		}
		fmt.Println("Hello was successful")
	}
}
