package channelz

import (
	"context"
	"fmt"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

func (h *grpcChannelzHandler) WriteSocketPage(w io.Writer, socket int64) {
	writeHeader(w, fmt.Sprintf("ChannelZ socket %d", socket))
	h.writeSocket(w, socket)
	writeFooter(w)
}

// writeSocket writes HTML to w containing RPC single socket stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *grpcChannelzHandler) writeSocket(w io.Writer, socket int64) {
	if err := socketTemplate.Execute(w, h.getSocket(socket)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *grpcChannelzHandler) getSocket(socketID int64) *channelzgrpc.GetSocketResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	socket, err := client.GetSocket(ctx, &channelzgrpc.GetSocketRequest{SocketId: socketID})
	if err != nil {
		log.Errorf("Error querying GetSocket %+v", err)
		return nil
	}
	return socket
}

const socketTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr classs="header">
        <th>Socket</th>
		<th>StreamsStarted</th>
        <th>StreamsSucceeded</th>
        <th>StreamsFailed</th>
        <th>MessagesSent</th>
        <th>MessagesReceived</th>
		<th>KeepAlivesSent</th>
		<th>LastLocalStreamCreated</th>
		<th>LastRemoteStreamCreated</th>
		<th>LastMessageSent</th>
		<th>LastMessageReceived</th>
		<th>LocalFlowControlWindow</th>
		<th>RemoteFlowControlWindow</th>
		<th>Options</th>
		<th>Security</th>
    </tr>
{{with .Socket}}
    <tr>
        <td>
			<b>{{.Ref.SocketId}}</b> {{.Ref.Name}}<br/>
			<!--
			<pre>{{.Local}} -> {{.Remote}} {{with .RemoteName}}({{.}}){{end}}</pre>
			-->
		</td>
		{{with .Data}}
			<td>{{.StreamsStarted}}</td>
			<td>{{.StreamsSucceeded}}</td>
			<td>{{.StreamsFailed}}</td>
			<td>{{.MessagesSent}}</td>
			<td>{{.MessagesReceived}}</td>
			<td>{{.KeepAlivesSent}}</td>
			<td>{{.LastLocalStreamCreatedTimestamp | timestamp}}</td>
			<td>{{.LastRemoteStreamCreatedTimestamp | timestamp}}</td>
			<td>{{.LastMessageSentTimestamp | timestamp}}</td>
			<td>{{.LastMessageReceivedTimestamp | timestamp}}</td>
			<td>{{.LocalFlowControlWindow.Value}}</td>
			<td>{{.RemoteFlowControlWindow.Value}}</td>
			<td>
				{{range .Option}}
					{{.Name}}: {{.Value}} {{with .Additional}}({{.}}){{end}}<br/>
				{{end}}
			</td>
		{{end}}
		<td>{{.Security}}</td>
	</tr>
{{end}}
</table>
`
