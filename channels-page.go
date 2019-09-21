package channelz

import (
	"context"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// WriteChannelsPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *grpcChannelzHandler) WriteChannelsPage(w io.Writer, start int64) {
	writeHeader(w, "Channels")
	h.writeChannels(w, start)
	writeFooter(w)
}

// writeTopChannels writes HTML to w containing per-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *grpcChannelzHandler) writeChannels(w io.Writer, start int64) {
	if err := channelsTemplate.Execute(w, h.getTopChannels(start)); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *grpcChannelzHandler) getTopChannels(start int64) *channelzgrpc.GetTopChannelsResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	channels, err := client.GetTopChannels(ctx, &channelzgrpc.GetTopChannelsRequest{
		StartChannelId: start,
	})
	if err != nil {
		log.Errorf("Error querying GetTopChannels %+v", err)
		return nil
	}
	return channels
}

const channelsTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr class="header">
		<th colspan=100 style="text-align:left">Top Channels: {{.Channel | len}}</th>
    </tr>

	{{template "channel-header"}}
	{{$last := .Channel}}
	{{range .Channel}}
		{{template "channel-body" .}}
		{{$last = .}}
	{{end}}
	{{if not .End}}
		<tr>
			<th colspan=100 style="text-align:left">
				<a href="{{link "channels"}}?start={{$last.Ref.ChannelId}}">Next&nbsp;&gt;</a>
			</th>
		</tr>
	{{end}}
</table>
<br/>
<br/>
`
