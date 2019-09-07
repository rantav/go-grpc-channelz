package channelz

import (
	"context"
	"fmt"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// WriteChannelPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *grpcChannelzHandler) WriteChannelPage(w io.Writer, channel int64) {
	writeHeader(w, fmt.Sprintf("ChannelZ channel %d", channel))
	h.writeChannel(w, channel)
	writeFooter(w)
}

func (h *grpcChannelzHandler) writeChannel(w io.Writer, channel int64) {
	if err := channelTemplate.Execute(w, h.getChannel(channel)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *grpcChannelzHandler) getChannel(channelID int64) *channelzgrpc.GetChannelResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	channel, err := client.GetChannel(ctx, &channelzgrpc.GetChannelRequest{ChannelId: channelID})
	if err != nil {
		log.Errorf("Error querying GetChannel %+v", err)
		return nil
	}
	return channel
}

const channelTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr classs="header">
        <th>Channel</th>
        <th>State</th>
        <th>Target</th>
        <th>Subchannels</th>
        <th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
        <th>ChannelRef</th>
    </tr>
{{with .Channel}}
    <tr>
        <td><b>{{.Ref.ChannelId}}</b> {{.Ref.Name}}</td>
        <td>{{.Data.State}}</td>
        <td>{{.Data.Target}}</td>
		<td>
			{{range .SubchannelRef}}
				<a href="{{link "subchannel" .SubchannelId}}"><b>{{.SubchannelId}}</b> {{.Name}}</a><br/>
			{{end}}
		</td>
        <td>{{.Data.Trace.CreationTimestamp | timestamp}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>{{.ChannelRef}}</td>
	</tr>
    <tr classs="header">
        <th colspan=100>Events</th>
    </tr>
	<tr>
		<td>&nbsp;</td>
        <td colspan=100>
			<pre>
			{{- range .Data.Trace.Events}}
{{.Severity}} [{{.Timestamp | timestamp}}]: {{.Description}}
			{{- end -}}
			</pre>
		</td>
    </tr>
{{end}}
</table>
`
