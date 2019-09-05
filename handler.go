package channelz

import (
	"net/http"
	"path"
	"sync"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
)

// CreateHandler creates an http handler with the routes of channelz mounted to the provided prefix.
// pathPrefix is the prefix to which /channelz will be prepended
// grpcBindAddress is the TCP bind address for the gRPC service you'd like to monitor.
// 	grpcBindAddress is required since the channelz interface connects to this gRPC service
//
// Typically you'd use the return value of CreateHandler as an argument to http.Handle
// For example:
// 	http.Handle("/", channelz.CreateHandler("/foo", grpcBindAddress))
func CreateHandler(pathPrefix, grpcBindAddress string) http.Handler {
	prefix := path.Join(pathPrefix, "channelz") + "/"
	handler := &channelzHandler{bindAddress: grpcBindAddress}
	return createRouter(prefix, handler)
}

type channelzHandler struct {
	// the target server's bind address
	bindAddress string

	// The client connection (lazily initialized)
	client channelzgrpc.ChannelzClient

	mu sync.Mutex
}
