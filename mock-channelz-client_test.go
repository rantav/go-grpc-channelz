package channelz

import (
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	channelzclient "google.golang.org/grpc/channelz/grpc_channelz_v1"
)

type mockChannelzClient struct{}

func (m *mockChannelzClient) GetTopChannels(
	ctx context.Context,
	in *channelzclient.GetTopChannelsRequest,
	opts ...grpc.CallOption) (*channelzclient.GetTopChannelsResponse, error) {
	return &channelzclient.GetTopChannelsResponse{
		Channel: []*channelzclient.Channel{
			createMockChannel(),
		},
	}, nil
}

func (m *mockChannelzClient) GetServers(
	ctx context.Context,
	in *channelzclient.GetServersRequest,
	opts ...grpc.CallOption) (*channelzclient.GetServersResponse, error) {
	return &channelzclient.GetServersResponse{
		Server: []*channelzclient.Server{
			createMockServer(),
		},
	}, nil
}

func createMockServer() *channelzclient.Server {
	return &channelzclient.Server{
		Ref: &channelzclient.ServerRef{
			ServerId: 1,
			Name:     "one",
		},
		Data: &channelzclient.ServerData{
			Trace:          createMockChannelTrace(),
			CallsStarted:   1,
			CallsSucceeded: 1,
			CallsFailed:    0,
			LastCallStartedTimestamp: &timestamp.Timestamp{
				Seconds: 6,
				Nanos:   7,
			},
		},
		ListenSocket: []*channelzclient.SocketRef{{
			SocketId: 6,
			Name:     "six",
		}},
	}
}

func (m *mockChannelzClient) GetServer(
	ctx context.Context,
	in *channelzclient.GetServerRequest,
	opts ...grpc.CallOption) (*channelzclient.GetServerResponse, error) {
	return &channelzclient.GetServerResponse{
		Server: createMockServer(),
	}, nil
}

func (m *mockChannelzClient) GetServerSockets(
	ctx context.Context,
	in *channelzclient.GetServerSocketsRequest,
	opts ...grpc.CallOption) (*channelzclient.GetServerSocketsResponse, error) {
	return nil, nil
}

func (m *mockChannelzClient) GetChannel(
	ctx context.Context,
	in *channelzclient.GetChannelRequest,
	opts ...grpc.CallOption) (*channelzclient.GetChannelResponse, error) {
	return &channelzclient.GetChannelResponse{
		Channel: createMockChannel(),
	}, nil
}

func createMockChannelTrace() *channelzclient.ChannelTrace {
	return &channelzclient.ChannelTrace{
		NumEventsLogged: 5,
		CreationTimestamp: &timestamp.Timestamp{
			Seconds: 6,
			Nanos:   7,
		},
		Events: []*channelzclient.ChannelTraceEvent{{
			Description: "setup",
			Severity:    channelzclient.ChannelTraceEvent_CT_INFO,
			Timestamp: &timestamp.Timestamp{
				Seconds: 6,
				Nanos:   7,
			},
		}},
	}
}
func createMockChannelData() *channelzclient.ChannelData {
	return &channelzclient.ChannelData{
		State: &channelzclient.ChannelConnectivityState{
			State: channelzclient.ChannelConnectivityState_CONNECTING,
		},
		Target:         "the world",
		Trace:          createMockChannelTrace(),
		CallsStarted:   1,
		CallsSucceeded: 2,
		CallsFailed:    0,
		LastCallStartedTimestamp: &timestamp.Timestamp{
			Seconds: 6,
			Nanos:   7,
		},
	}
}
func createMockChannel() *channelzclient.Channel {
	return &channelzclient.Channel{
		Ref: &channelzclient.ChannelRef{
			ChannelId: 5,
			Name:      "five",
		},
		Data: createMockChannelData(),
		ChannelRef: []*channelzclient.ChannelRef{{
			ChannelId: 7,
			Name:      "seven",
		}},
		SubchannelRef: []*channelzclient.SubchannelRef{{
			SubchannelId: 8,
			Name:         "eight",
		}},
	}
}

func (m *mockChannelzClient) GetSubchannel(
	ctx context.Context, in *channelzclient.GetSubchannelRequest,
	opts ...grpc.CallOption) (*channelzclient.GetSubchannelResponse, error) {
	return &channelzclient.GetSubchannelResponse{
		Subchannel: createMockSubchannel(),
	}, nil
}

func createMockSubchannel() *channelzclient.Subchannel {
	return &channelzclient.Subchannel{
		Ref: &channelzclient.SubchannelRef{
			SubchannelId: 4,
			Name:         "four",
		},
		Data:          createMockChannelData(),
		ChannelRef:    []*channelzclient.ChannelRef{},
		SubchannelRef: []*channelzclient.SubchannelRef{},
		SocketRef: []*channelzclient.SocketRef{{
			SocketId: 9,
			Name:     "nine",
		}},
	}
}

func (m *mockChannelzClient) GetSocket(
	ctx context.Context,
	in *channelzclient.GetSocketRequest,
	opts ...grpc.CallOption) (*channelzclient.GetSocketResponse, error) {
	return nil, nil
}
