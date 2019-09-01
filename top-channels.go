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
	writeFooter(w)
}

func writeHeader(w io.Writer, rtitle string) {
	if err := headerTemplate.Execute(w, headerData{Title: "ChannelZ Stats"}); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func writeFooter(w io.Writer) {
	if err := footerTemplate.Execute(w, nil); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}

}

// writeHTMLChannelzSummary writes HTML to w containing per-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *channelzHandler) writeTopChannels(w io.Writer) {
	if err := channelzsTemplate.Execute(w, h.getTopChannels()); err != nil {
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
