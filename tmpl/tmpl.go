package tmpl

import (
	"embed"
	"mime"
	"path"
	"text/template"
	"time"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

func init() {
	tmpl = template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"datef": func(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) },
		"mime":  func(s string) string { return mime.TypeByExtension(path.Ext(s)) },
	})
	tmpl = template.Must(tmpl.ParseFS(tmplFS, "templates/*"))
}
