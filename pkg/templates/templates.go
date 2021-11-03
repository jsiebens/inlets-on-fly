package templates

import (
	"embed"
	"io"
	"text/template"
)

//go:embed *.tpl
var fs embed.FS

var t *template.Template

func init() {
	t = template.Must(template.ParseFS(fs, "*.tpl"))
}

func Render(w io.Writer, name string, data interface{}) error {
	return t.ExecuteTemplate(w, name, data)
}
