package channelz

import (
	"context"
	"fmt"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// writeTopChannelsPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *channelzHandler) writeSubchannelPage(w io.Writer, subchannel int64) {
	writeHeader(w, fmt.Sprintf("ChannelZ subchannel %d", subchannel))
	h.writeSubchannel(w, subchannel)
	writeFooter(w)
}

// writeSubchannel writes HTML to w containing sub-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *channelzHandler) writeSubchannel(w io.Writer, subchannel int64) {
	if err := subChannelTemplate.Execute(w, h.getSubchannel(subchannel)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *channelzHandler) getSubchannel(subchannelID int64) *channelzgrpc.GetSubchannelResponse {
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
<table frame=box cellspacing=0 cellpadding=2>
    <tr classs="header">
        <th>Subchannel</th>
        <th>State</th>
        <th>Target</th>
        <th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
        <th>ChannelRef</th>
        <th>SocketRef</th>

    </tr>
{{with .Subchannel}}
    <tr>
        <td><b>{{.Ref.SubchannelId}}</b> {{.Ref.Name}}</td>
        <td>{{.Data.State}}</td>
        <td>{{.Data.Target}}</td>
        <td>{{.Data.Trace.CreationTimestamp | timestamp}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>{{.ChannelRef}}</td>
		<td>
			{{range .SocketRef}}
				<b>{{.SocketId}}</b> {{.Name}}<br/>
			{{end}}
		</td>
	</tr>
	<tr>
        <th colspan=100>Events</th>
	</tr>
	<tr>
        <td colspan=100>
			<pre>
			{{- range .Data.Trace.Events}}
{{.Severity}} [{{.Timestamp}}]: {{.Description}}
			{{- end -}}
			</pre>
		</td>
    </tr>
{{end}}
</table>
`
