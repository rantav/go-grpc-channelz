package channelz

// channelz is a gRPC protocol for exposing channels usage.
// https://grpc.io/blog/a_short_introduction_to_channelz/
// Channels in gRPC are simply connections, so channelz is the exposition
// of clients or server connections.
// This library exposes channelz by making an RPC to the localhost.
// For this to work the channelz service should also be registered with
// "google.golang.org/grpc/channelz/service".RegisterChannelzServiceToServer(grpcServer)

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
