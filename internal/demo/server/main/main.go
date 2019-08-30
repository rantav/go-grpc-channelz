package main

import (
	"fmt"
	"net"

	log "google.golang.org/grpc/grpclog"

	"github.com/rantav/go-grpc-channelz/internal/demo/server"
)

func main() {
	bindAddress := ":8080"
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer, err := server.New()
	if err != nil {
		log.Fatalf("Failed to create grpc server %+v", err)
	}

	fmt.Printf("demo server is up is up; bind address: %s \n", bindAddress)

	// should never return
	err = grpcServer.Serve(listener)

	log.Fatalf("Error running server: %+v", err)

}
