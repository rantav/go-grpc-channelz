package channelz

import (
	"context"
	"fmt"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// WriteSubchannelsPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *grpcChannelzHandler) WriteSubchannelPage(w io.Writer, subchannel int64) {
	writeHeader(w, fmt.Sprintf("ChannelZ subchannel %d", subchannel))
	h.writeSubchannel(w, subchannel)
	writeFooter(w)
}

// writeSubchannel writes HTML to w containing sub-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *grpcChannelzHandler) writeSubchannel(w io.Writer, subchannel int64) {
	if err := subChannelTemplate.Execute(w, h.getSubchannel(subchannel)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *grpcChannelzHandler) getSubchannel(subchannelID int64) *channelzgrpc.GetSubchannelResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	subchannel, err := client.GetSubchannel(ctx, &channelzgrpc.GetSubchannelRequest{
		SubchannelId: subchannelID,
	})
	if err != nil {
		log.Errorf("Error querying GetSubchannel %+v", err)
		return nil
	}
	return subchannel
}

const subChannelsTemplateHTML = `
{{define "subchannel-header"}}
    <tr classs="header">
        <th>Subchannel</th>
        <th>State</th>
        <th>Target</th>
        <th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
        <th>Child Channels</th>
        <th>Child Subchannels</th>
        <th>Socket</th>
    </tr>
{{end}}

{{define "subchannel-body"}}
    <tr>
        <td><a href="{{link "subchannel" .Ref.SubchannelId}}"><b>{{.Ref.SubchannelId}}</b> {{.Ref.Name}}</a></td>
        <td>{{.Data.State}}</td>
        <td>{{.Data.Target}}</td>
        <td>{{.Data.Trace.CreationTimestamp | timestamp}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>
			{{range .ChannelRef}}
				<b><a href="{{link "channel" .ChannelId}}">{{.ChannelId}}</b> {{.Name}}</a><br/>
			{{end}}
		</td>
		<td>
			{{range .SubchannelRef}}
				<b><a href="{{link "subchannel" .SubchannelId}}">{{.SubchannelId}}</b> {{.Name}}</a><br/>
			{{end}}
		</td>
		<td>
			{{range .SocketRef}}
				<b><a href="{{link "socket" .SocketId}}">{{.SocketId}}</b> {{.Name}}</a><br/>
			{{end}}
		</td>
	</tr>
{{end}}

{{define "subchannel-events"}}
	<tr>
        <th colspan=100>Events</th>
	</tr>
	<tr>
        <td colspan=100>
			<pre>
			{{- range .Data.Trace.Events}}
{{.Severity}} [{{.Timestamp | timestamp}}]: {{.Description}}
			{{- end -}}
			</pre>
		</td>
    </tr>
{{end}}

<table frame=box cellspacing=0 cellpadding=2>
	{{template "subchannel-header"}}
	{{template "subchannel-body" .Subchannel}}
	{{template "subchannel-events" .Subchannel}}
</table>
`
