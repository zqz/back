package app

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/pressly/chi/render"
)

// The CSS Assets to inject into the page.
func templateData(host string) map[string]interface{} {
	css := []string{
		"login", "menu", "registration", "table", "authentication",
		"uploader", "upload_file", "style", "alerts", "dashboard",
		"header", "footer", "file_list_component", "file_view",
	}
	js := []string{
		"new/xhr", "ws",
		"lib/dominate", "lib/filesize", "helpers", "alerts", "table_component",
		"authentication_component", "router", "login", "login_component", "menu",
		"header_component", "footer_component", "app",
	}

	nakedHost := strings.Split(host, ":")[0]
	var liveReloadStr string
	if config.LiveReload {
		liveReloadStr = "http://" + nakedHost + ":35729/livereload.js?snipver=1"
	}

	var wsPath string
	if config.Secure {
		wsPath = fmt.Sprintf("wss://%s/ws", host)
	} else {
		wsPath = fmt.Sprintf("ws://%s/ws", host)
	}

	return map[string]interface{}{
		"Title":      "zqz.ca",
		"LiveReload": liveReloadStr,
		"Cdn":        template.JSStr(config.CDNURL),
		"WSPath":     template.JSStr(wsPath),

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
      <meta name="viewport" content="width=device-width, user-scalable=no">
      <link rel='shortcut icon' href='{{ .Cdn }}/favicon.ico'/>
      {{- range .Assets.Css }}
      <link rel='stylesheet' media='screen' href='{{ $cdn }}/{{ . }}.css'/>
      {{- end }}
    </head>
    <body>
      <script type='text/javascript'>
        window.cdn = {{$cdn}};
        window.ws_url = {{.WSPath}};
      </script>
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

// Index generates an index.html
func Index(w http.ResponseWriter, r *http.Request) {
	d := templateData(r.Host)
	o := generateIndex(d)
	render.HTML(w, r, o)
}
