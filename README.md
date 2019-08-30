# go-grpc-channelz
A ChannelZ UI for gRPC in Golang

## Usage


```go
import (
	channelzservice "google.golang.org/grpc/channelz/service"
)

.
.
.

// Register the channelz handler
channelz.Handle(http.DefaultServeMux, "/")

// Register the channelz gRPC service to grpcServer so that we can query it for this service.
channelzservice.RegisterChannelzServiceToServer(grpcServer)

// Listen and serve HTTP for the default serve mux
adminListener, err := net.Listen("tcp", ":8081")
if err != nil {
    log.Fatal(err)
}
go http.Serve(adminListener, nil)

```
