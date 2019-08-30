package channelz

// channelz is a gRPC protocol for exposing channels usage.
// https://grpc.io/blog/a_short_introduction_to_channelz/
// Channels in gRPC are simply connections, so channelz is the exposition
// of clients or server connections.
// This library exposes channelz by making an RPC to the localhost.
// For this to work the channelz service should also be registered with
// "google.golang.org/grpc/channelz/service".RegisterChannelzServiceToServer(grpcServer)

import (
	"io"
	"net/http"
	"path"

	log "google.golang.org/grpc/grpclog"
)

// Handle adds the /channelz to the given ServeMux rooted at pathPrefix.
// if mus is nill then http.DefaultServeMux is used.
func Handle(mux *http.ServeMux, pathPrefix string) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	mux.HandleFunc(path.Join(pathPrefix, "channelz"), channelzHandler)
	// mux.HandleFunc(path.Join(pathPrefix, "tracez"), tracezHandler)
	// mux.Handle(path.Join(pathPrefix, "public/"), http.FileServer(fs))
}

func channelzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	WriteHTMLChannelzPage(w)
}

// WriteHTMLChannelzPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func WriteHTMLChannelzPage(w io.Writer) {
	if err := headerTemplate.Execute(w, headerData{Title: "ChannelZ Stats"}); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
	WriteHTMLChannelzSummary(w)
	if err := footerTemplate.Execute(w, nil); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

// WriteHTMLChannelzSummary writes HTML to w containing per-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func WriteHTMLChannelzSummary(w io.Writer) {
	if err := channelzsTemplate.Execute(w, getChannelzsList()); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func getChannelzsList() string {
	return "TODO"
}

// func NewChannelzClient(dialString string) (channelzgrpc.ChannelzClient, error) {
// 	conn, err := grpc.Dial(dialString, grpc.WithInsecure())
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "error dialing to %s", dialString)
// 	}
// 	client := channelzgrpc.NewChannelzClient(conn)
// 	return client, nil
// }

// func MonitorChannelz(bindAddress string) {
// 	host := getHostFromBindAddress(bindAddress)
// 	channelzclient, err := NewChannelzClient(host)
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

// func getHostFromBindAddress(bindAddress string) string {
// 	if strings.HasPrefix(bindAddress, ":") {
// 		return fmt.Sprintf("localhost%s", bindAddress)
// 	}
// 	return bindAddress
// }
