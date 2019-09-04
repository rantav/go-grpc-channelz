package channelz

// channelz is a gRPC protocol for exposing channels usage.
// https://grpc.io/blog/a_short_introduction_to_channelz/
// Channels in gRPC are simply connections, so channelz is the exposition
// of clients or server connections.
// This library exposes channelz by making an RPC to the localhost.
// For this to work the channelz service should also be registered with
// "google.golang.org/grpc/channelz/service".RegisterChannelzServiceToServer(grpcServer)
//
//
//

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
