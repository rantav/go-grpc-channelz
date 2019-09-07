package channelz

import (
	"context"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// writeServers writes HTML to w containing RPC servers stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *grpcChannelzHandler) writeServers(w io.Writer) {
	if err := serversTemplate.Execute(w, h.getServers()); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *grpcChannelzHandler) getServers() *channelzgrpc.GetServersResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	servers, err := client.GetServers(ctx, &channelzgrpc.GetServersRequest{})
	if err != nil {
		log.Errorf("Error querying GetServers %+v", err)
		return nil
	}
	return servers
}

const serversTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr class="header">
		<th colspan=100 style="text-align:left">Servers: {{.Server | len}}</th>
    </tr>

	{{template "server-header"}}
	{{range .Server}}
		{{template "server-body" .}}
	{{end}}
</table>
`
