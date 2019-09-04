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

var channelRx = regexp.MustCompile(`channelz/channel/(?P<channel>\d+)$`)
var subchannelRx = regexp.MustCompile(`channelz/subchannel/(?P<subchannel>\d+)$`)
var serverRx = regexp.MustCompile(`channelz/server/(?P<server>\d+)$`)

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
	if match := channelRx.FindStringSubmatch(path); match != nil {
		channel, err := strconv.ParseInt(match[1], 10, 0)
		if err != nil {
			log.Errorf("channelz: Unable to parse int for channel ID. %s", match[1])
		}
		h.writeChannelPage(w, channel)
		return
	}
	if match := subchannelRx.FindStringSubmatch(path); match != nil {
		subChannel, err := strconv.ParseInt(match[1], 10, 0)
		if err != nil {
			log.Errorf("channelz: Unable to parse int for sub-channel ID. %s", match[1])
		}
		h.writeSubchannelPage(w, subChannel)
		return
	}
	if match := serverRx.FindStringSubmatch(path); match != nil {
		server, err := strconv.ParseInt(match[1], 10, 0)
		if err != nil {
			log.Errorf("channelz: Unable to parse int for server ID. %s", match[1])
		}
		h.writeServerPage(w, server)
		return
	}
	write404(w)
}
