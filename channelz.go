package channelz

// channelz is a gRPC protocol for exposing channels usage.
// https://grpc.io/blog/a_short_introduction_to_channelz/
// Channels in gRPC are simply connections, so channelz is the exposition
// of clients or server connections.
// This library exposes channelz by making an RPC to the localhost.
// For this to work the channelz service should also be registered with
// "google.golang.org/grpc/channelz/service".RegisterChannelzServiceToServer(grpcServer)

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
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
	h.writeHTMLChannelzPage(w)
}

// writeHTMLChannelzPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *channelzHandler) writeHTMLChannelzPage(w io.Writer) {
	if err := headerTemplate.Execute(w, headerData{Title: "ChannelZ Stats"}); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
	h.writeHTMLChannelzSummary(w)
	if err := footerTemplate.Execute(w, nil); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

// writeHTMLChannelzSummary writes HTML to w containing per-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *channelzHandler) writeHTMLChannelzSummary(w io.Writer) {
	if err := channelzsTemplate.Execute(w, h.getChannelzsList()); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *channelzHandler) getChannelzsList() *channelzgrpc.GetTopChannelsResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	channels, err := client.GetTopChannels(ctx, &channelzgrpc.GetTopChannelsRequest{})
	if err != nil {
		log.Errorf("Error querying GetTopChannels %+v", err)
		return nil
	}
	return channels
}

func (h *channelzHandler) connect() (channelzgrpc.ChannelzClient, error) {
	if h.client != nil {
		// Already connected
		return h.client, nil
	}

	host := getHostFromBindAddress(h.bindAddress)
	h.mu.Lock()
	defer h.mu.Unlock()
	client, err := newChannelzClient(host)
	if err != nil {
		return nil, err
	}
	h.client = client
	return h.client, nil
}

func newChannelzClient(dialString string) (channelzgrpc.ChannelzClient, error) {
	conn, err := grpc.Dial(dialString, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "error dialing to %s", dialString)
	}
	client := channelzgrpc.NewChannelzClient(conn)
	return client, nil
}

// func MonitorChannelz(bindAddress string) {
// 	host := getHostFromBindAddress(bindAddress)
// 	channelzclient, err := newChannelzClient(host)
// 	if err != nil {
// 		log.GetDefault().Errorf("Error creating channelz client %+v", err)
// 	}
// 	ctx := context.Background()
// 	for {
// 		time.Sleep(10 * time.Second)
// 		channels, err := channelzclient.GetTopChannels(ctx, &channelzgrpc.GetTopChannelsRequest{})
// 		if err != nil {
// 			log.GetDefault().Errorf("Error querying GetTopChannels %+v", err)
// 			continue
// 		}
// 		log.GetDefault().Infof("gRPC number of channels returned from GetTopChannels: %d", len(channels.Channel))
// 		for _, c := range channels.Channel {
// 			log.GetDefault().Debugf("\tChannel %d: %+v", c.GetRef().ChannelId, c)
// 			for _, scref := range c.SubchannelRef {
// 				subchannel, err := channelzclient.GetSubchannel(ctx, &channelzgrpc.GetSubchannelRequest{
// 					SubchannelId: scref.SubchannelId,
// 				})
// 				if err != nil {
// 					log.GetDefault().Errorf("Error querying GetSubchannel %+v", err)
// 					continue
// 				}
// 				log.GetDefault().Debugf("\t\t Channel %d, Subchannel %d: %+v",
// 					c.GetRef().ChannelId,
// 					subchannel.GetSubchannel().GetRef().GetSubchannelId(),
// 					subchannel)
// 			}
// 		}

// 		servers, err := channelzclient.GetServers(ctx, &channelzgrpc.GetServersRequest{})
// 		if err != nil {
// 			log.GetDefault().Errorf("Error querying GetServers %+v", err)
// 			continue
// 		}
// 		log.GetDefault().Infof("gRPC number of servers returned from GetServers: %d", len(servers.Server))
// 		for _, s := range servers.Server {
// 			log.GetDefault().Debugf("\tServer %d: %+v", s.GetRef().GetServerId(), s)
// 		}
// 	}
// }

func getHostFromBindAddress(bindAddress string) string {
	if strings.HasPrefix(bindAddress, ":") {
		return fmt.Sprintf("localhost%s", bindAddress)
	}
	return bindAddress
}
