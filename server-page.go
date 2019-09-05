package channelz

import (
	"context"
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

func (h *channelzHandler) writeServerPage(w io.Writer, server int64) {
	writeHeader(w, fmt.Sprintf("ChannelZ server %d", server))
	h.writeServer(w, server)
	writeFooter(w)
}

// writeServer writes HTML to w containing RPC single server stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *channelzHandler) writeServer(w io.Writer, server int64) {
	if err := serverTemplate.Execute(w, h.getServer(server)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *channelzHandler) getServer(serverID int64) *channelzgrpc.GetServerResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	server, err := client.GetServer(ctx, &channelzgrpc.GetServerRequest{ServerId: serverID})
	if err != nil {
		log.Errorf("Error querying GetServers %+v", err)
		return nil
	}
	spew.Dump(server)
	return server
}

const serverTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr classs="header">
        <th>Server</th>
		<th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
		<th>Sockets</th>
    </tr>
{{with .Server}}
    <tr>
        <td><b>{{.Ref.ServerId}}</b> {{.Ref.Name}}</td>
        <td>{{with .Data.Trace}} {{.CreationTimestamp | timestamp}} {{end}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>
			{{range .ListenSocket}}
				<b>{{.SocketId}}</b> {{.Name}} <br/>
			{{end}}
		</td>
	</tr>
	{{with .Data.Trace}}
		<tr classs="header">
			<th colspan=100>Events</th>
		</tr>
		<tr>
			<td>&nbsp;</td>
			<td colspan=100>
				<pre>
				{{- range .Events}}
{{.Severity}} [{{.Timestamp}}]: {{.Description}}
				{{- end -}}
				</pre>
			</td>
		</tr>
	{{end}}
{{end}}
</table>
`