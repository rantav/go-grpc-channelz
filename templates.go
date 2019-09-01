package channelz

import (
	"text/template"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

var (
	headerTemplate      = parseTemplate("header", headerTemplateHTML)
	topChannelsTemplate = parseTemplate("channelzs", topChannelsTemplateHTML)
	subChannelTemplate  = parseTemplate("channelzs", subChannelsTemplateHTML)
	footerTemplate      = parseTemplate("footer", footerTemplateHTML)
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
	headerTemplateHTML = `
<!DOCTYPE html>
<html lang="en"><head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.indigo-pink.min.css">
	<style>
		body {padding: 1em}
		table {
			background-color: #fff5ee;
		}
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

	footerTemplateHTML = `
</body>
</html>
`
)
