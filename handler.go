package channelz

import (
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

var subchannelRx = regexp.MustCompile(`channelz/(?P<channel>\d+)/(?P<subchannel>\d+)$`)

// Handle adds the /channelz to the given ServeMux rooted at pathPrefix.
// if mux is nill then http.DefaultServeMux is used.
// pathPrefix is the prefix to which /channelz will be prepended
// bindAddress is the TCP bind address for the gRPC service you'd like to monitor.
// 	bindAddress is required since the channelz interface connects to this gRPC service
func Handle(mux *http.ServeMux, pathPrefix string, bindAddress string) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	mux.Handle(path.Join(pathPrefix, "channelz")+"/", &channelzHandler{
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
	path := r.URL.EscapedPath()
	if strings.HasSuffix(path, "channelz/") {
		h.writeTopChannelsPage(w)
		return
	}
	if match := subchannelRx.FindStringSubmatch(path); match != nil {
		topChannel, err := strconv.ParseInt(match[1], 10, 0)
		if err != nil {
			log.Errorf("channelz: Unable to parse int for channel ID. %s", match[1])
		}
		subChannel, err := strconv.ParseInt(match[2], 10, 0)
		if err != nil {
			log.Errorf("channelz: Unable to parse int for sub-channel ID. %s", match[2])
		}
		h.writeSubchannelPage(w, topChannel, subChannel)
		return
	}
	write404(w)
}
