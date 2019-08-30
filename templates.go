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

var (
	channelzsTemplateHTML = `
<h1>HELLO WORLD</h1>
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
<body>
<h1>{{.Title}}</h1>
`
)
