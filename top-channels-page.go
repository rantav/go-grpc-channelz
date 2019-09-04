package channelz

import (
	"context"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// writeTopChannelsPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *channelzHandler) writeTopChannelsPage(w io.Writer) {
	writeHeader(w, "ChannelZ Stats")
	h.writeTopChannels(w)
	h.writeServers(w)
	writeFooter(w)
}

func writeHeader(w io.Writer, title string) {
	if err := headerTemplate.Execute(w, headerData{Title: title}); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func writeFooter(w io.Writer) {
	if err := footerTemplate.Execute(w, nil); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}

}

// writeTopChannels writes HTML to w containing per-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *channelzHandler) writeTopChannels(w io.Writer) {
	if err := topChannelsTemplate.Execute(w, h.getTopChannels()); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *channelzHandler) getTopChannels() *channelzgrpc.GetTopChannelsResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	channels, err := client.GetTopChannels(ctx, &channelzgrpc.GetTopChannelsRequest{})
	if err != nil {
		log.Errorf("Error querying GetTopChannels %+v", err)
		return nil
	}
	return channels
}

const topChannelsTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr class="header">
		<th colspan=100 style="text-align:left">Top Channels: {{.Channel | len}}</th>
    </tr>

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
{{range .Channel}}
    <tr>
        <td>
			<a href="channel/{{.Ref.ChannelId}}">[{{.Ref.ChannelId}}] {{.Ref.Name}}</a>
		</td>
        <td>{{.Data.State}}</td>
        <td>{{.Data.Target}}</td>
		<td>
			{{range .SubchannelRef}}
				<a href="subchannel/{{.SubchannelId}}">[{{.SubchannelId}}] {{.Name}}</a><br/>
			{{end}}
		</td>
        <td>{{.Data.Trace.CreationTimestamp | timestamp}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>{{.ChannelRef}}</td>
	</tr>
{{end}}
</table>
`
