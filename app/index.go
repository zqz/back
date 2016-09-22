package app

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/zqzca/echo"
)

// The CSS Assets to inject into the page.
func templateData(host string) map[string]interface{} {
	css := []string{
		"login", "menu", "registration", "table", "authentication",
		"uploader", "upload_file", "style", "alerts", "dashboard",
		"header", "footer", "file_list_component", "file_view",
	}
	js := []string{
		"new/xhr",
		"lib/dominate", "lib/filesize", "helpers", "alerts", "table_component",
		"authentication_component", "router", "login", "login_component", "menu",
		"header_component", "footer_component", "app",
	}

	var liveReloadStr string
	if config.LiveReload {
		liveReloadStr = "http://" + host + ":35729/livereload.js?snipver=1"
	}

	return map[string]interface{}{
		"Title":      "zqz.ca",
		"LiveReload": liveReloadStr,
		"Cdn":        template.JSStr(config.CDNURL),

		"Assets": map[string]interface{}{
			"Js":  js,
			"Css": css,
		},
	}
}

func generateIndex(tmplData map[string]interface{}) string {
	indexTemplate := `{{ $cdn := .Cdn -}}
<!DOCTYPE HTML>
  <html>
    <head>
      <meta http-equiv='content-type' content='text/html; charset=utf-8'>
      <title>zqz.ca</title>
      <link rel='shortcut icon' href='{{ .Cdn }}/favicon.ico'/>
      {{- range .Assets.Css }}
      <link rel='stylesheet' media='screen' href='{{ $cdn }}/{{ . }}.css'/>
      {{- end }}
    </head>
    <body>
      <script type='text/javascript'>window.cdn = {{$cdn}};</script>
      {{- with .LiveReload }}
      <script type='text/javascript' src='{{.}}'></script>
      {{- end }}
      {{- range .Assets.Js }}
      <script type='text/javascript' src='{{$cdn}}/{{.}}.js'></script>
      {{- end }}
    </body>
  </html>`

	t := template.New("App Index Template")
	t, err := t.Parse(indexTemplate)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	err = t.Execute(&output, tmplData)
	if err != nil {
		panic(err)
	}

	return output.String()
}

// AppIndex generates an index.html
func AppIndex(c echo.Context) error {
	host := strings.Split(c.Request().Host(), ":")[0]
	d := templateData(host)
	o := generateIndex(d)
	return c.HTML(200, o)
}
