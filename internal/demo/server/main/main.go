package main

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
	log "google.golang.org/grpc/grpclog"

	channelz "github.com/rantav/go-grpc-channelz"
	"github.com/rantav/go-grpc-channelz/internal/demo/server"
)

func main() {
	const (
		grpcBindAddress  = ":8080"
		adminBindAddress = ":8081"
	)

	grpcListener, err := net.Listen("tcp", grpcBindAddress)
	if err != nil {
		log.Fatal(err)
	}

	adminListener, err := net.Listen("tcp", adminBindAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer, err := server.New()
	if err != nil {
		log.Fatalf("Failed to create grpc server %+v", err)
	}

	// Register the channelz handler
	channelz.Handle(http.DefaultServeMux, "/")

	g := new(errgroup.Group)
	g.Go(func() error { return http.Serve(adminListener, nil) })
	g.Go(func() error { return grpcServer.Serve(grpcListener) })

	fmt.Printf("demo server is up is up; gRPC bind address: %s, http admin address: %s \n", grpcBindAddress, adminBindAddress)

	// should never return
	err = g.Wait()
	log.Fatalf("Error running server: %+v", err)
}
