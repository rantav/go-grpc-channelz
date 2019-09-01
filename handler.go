package channelz

import (
	"net/http"
	"path"
	"sync"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
)

// Handle adds the /channelz to the given ServeMux rooted at pathPrefix.
// if mux is nill then http.DefaultServeMux is used.
// pathPrefix is the prefix to which /channelz will be prepended
// bindAddress is the TCP bind address for the gRPC service you'd like to monitor.
// 	bindAddress is required since the channelz interface connects to this gRPC service
func Handle(mux *http.ServeMux, pathPrefix string, bindAddress string) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	mux.Handle(path.Join(pathPrefix, "channelz"), &channelzHandler{
		bindAddress: bindAddress,
	})
	// mux.HandleFunc(path.Join(pathPrefix, "tracez"), tracezHandler)
	// mux.Handle(path.Join(pathPrefix, "public/"), http.FileServer(fs))
}

type channelzHandler struct {
	// the target server's bind address
	bindAddress string

	// The client connection (lazily initialized)
	client channelzgrpc.ChannelzClient

	mu sync.Mutex
}

func (h *channelzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.writeTopChannelsPage(w)
}
