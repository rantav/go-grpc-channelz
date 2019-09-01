package channelz

import (
	"text/template"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

var (
	headerTemplate    = parseTemplate("header", headerTemplateHTML)
	channelzsTemplate = parseTemplate("channelzs", channelzsTemplateHTML)
	footerTemplate    = parseTemplate("footer", footerTemplateHTML)
)

func parseTemplate(name, html string) *template.Template {
	return template.Must(template.New(name).Funcs(getFuncs()).Parse(html))
}

func getFuncs() template.FuncMap {
	return template.FuncMap{
		"timestamp": formatTimestamp,
	}
}

func formatTimestamp(ts *timestamp.Timestamp) string {
	t := time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	return t.Format(time.RFC3339)
}

// headerData contains data for the header template.
type headerData struct {
	Title string
}

var (
	channelzsTemplateHTML = `
<table bgcolor="#fff5ee" frame=box cellspacing=0 cellpadding=2>
    <tr class="header">
		<th colspan=100 style="text-align:left">Top Channels: {{.Channel | len}}</th>
    </tr>

    <tr classs="header">
        <th>ID</th>
        <th>Name</th>
        <th>State</th>
        <th>Target</th>
        <th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
        <th>ChannelRef</th>
        <th>SubchannelRef</th>
        <th>SocketRef</th>

        <th>Events</th>
    </tr>
{{range .Channel}}
    <tr>
        <td>{{.Ref.ChannelId}}</td>
        <td><b>{{.Ref.Name}}</b></td>
        <td>{{.Data.State}}</td>
        <td>{{.Data.Target}}</td>
        <td>{{.Data.Trace.CreationTimestamp | timestamp}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>{{.ChannelRef}}</td>
		<td>{{.SubchannelRef}}</td>
		<td>{{.SocketRef}}</td>
        <td>
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

	footerTemplateHTML = `
</body>
</html>
`

	headerTemplateHTML = `
<!DOCTYPE html>
<html lang="en"><head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.indigo-pink.min.css">
	<style>
		body {padding: 1em}
		tr.header {
			background-color: #eee5de;
		}
		td {
			vertical-align: top;
		}
	</style>
</head>
<body>
<h1>{{.Title}}</h1>
`
)
