package channelz

import (
	"text/template"
)

var (
	headerTemplate    = parseTemplate("header", headerTemplateHTML)
	channelzsTemplate = parseTemplate("channelzs", channelzsTemplateHTML)
	footerTemplate    = parseTemplate("footer", footerTemplateHTML)
)

func parseTemplate(name, html string) *template.Template {
	return template.Must(template.New(name).Parse(html))
}

// headerData contains data for the header template.
type headerData struct {
	Title string
}

type channels struct {
	Channels []channel
}

type channel struct {
}

var (
	channelzsTemplateHTML = `
<table bgcolor="#fff5ee" frame=box cellspacing=0 cellpadding=2>
    <tr bgcolor="#eee5de">
		<th class="l1" colspan=3>Top Channels: {{.Channel | len}}</th>
    </tr>
    <tr bgcolor="#eee5de">
        <th class="l1" colspan=3>Hello</th>
    </tr>
{{range .Channel}}
    <tr>
        <td><b>{{.Ref}}</b></td>
        <td></td>
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
    <script defer src="https://code.getmdl.io/1.3.0/material.min.js"></script>
</head>
<body style="padding: 2em">
<h1>{{.Title}}</h1>
`
)
