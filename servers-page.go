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
func (h *channelzHandler) writeServers(w io.Writer) {
	if err := serversTemplate.Execute(w, h.getServers()); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *channelzHandler) getServers() *channelzgrpc.GetServersResponse {
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

    <tr classs="header">
        <th>ID</th>
        <th>Name</th>
		<th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
    </tr>
{{range .Server}}
    <tr>
        <td>{{.Ref.ServerId}}</td>
        <td><b>{{.Ref.Name}}</b></td>
        <td>{{with .Data.Trace}} {{.CreationTimestamp | timestamp}} {{end}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
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
