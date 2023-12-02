package channelz

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/rantav/go-grpc-channelz/internal/demo/client"
	"github.com/rantav/go-grpc-channelz/internal/demo/server"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	channelzSrv "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Just a simple smoke test
func TestCreateHandler(t *testing.T) {
	h := CreateHandler("/prefix", ":8080")
	assert.NotNil(t, h)
}

// TestCreateHandlerWithDialOpts performs an end-to-end test with custom dialOpts.
// It is using an in memory grpc demo server which also acts as a channelz provider.
// A goroutine will perform demo calls to the server and the tests search the
// output html for hints that everything is working.
func TestCreateHandlerWithDialOpts(t *testing.T) {
	d, _ := t.Deadline()
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	errs, ctx := errgroup.WithContext(ctx)

	// in memory listener for demo grpc server
	listener := bufconn.Listen(1024 * 1024)

	// demo grpc server including channelz
	demoServer, err := server.New()
	require.Nil(t, err)
	channelzSrv.RegisterChannelzServiceToServer(demoServer)

	// serve demo grpc server
	errs.Go(func() error {
		return demoServer.Serve(listener)
	})

	// cleanup grpc server once ctx is done
	errs.Go(func() error {
		<-ctx.Done()
		demoServer.Stop()
		return nil
	})

	// demo client loop
	errs.Go(func() error {
		// give the grpc server some time to spin up
		time.Sleep(time.Millisecond * 100)

		// construct demo client
		demoClient, err := client.NewWithDialOpts("",
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
					return listener.DialContext(ctx)
				}),
			}...)
		require.Nil(t, err)

		// create demoClient.Hello calls each second until ctx is done
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(time.Second):
				_, err := demoClient.Hello(ctx, &emptypb.Empty{})
				if err != nil {
					return err
				}
			}
		}
	})

	// give the demo client a chance to do a request before performing tests
	time.Sleep(time.Second)

	// create http.Handler which uses dialOpts to dial to in memory listener
	handler := CreateHandlerWithDialOpts("/prefix", ":8080",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return listener.DialContext(ctx)
			}),
		}...)

	// start a testing http server with handler
	s := httptest.NewServer(handler)
	defer s.Close()

	// actual test: get main page
	resp, err := http.Get(s.URL + "/prefix/channelz/")
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	body := string(bodyBytes)

	assert.Contains(t, string(body), "Servers: 1")
	assert.Contains(t, string(body), "/prefix/channelz/server/")
	assert.Contains(t, string(body), "state:READY")

	// teardown and cleanup
	cancel()
	err = errs.Wait()
	require.Nil(t, err)
}
