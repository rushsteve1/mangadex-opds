package tmpl

import (
	"embed"
	"mime"
	"path"
	"regexp"
	"text/template"
	"time"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

func init() {
	tmpl = template.New("")
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	tmpl = tmpl.Funcs(template.FuncMap{
		"add":    func(a int, b int) int { return a + b },
		"id":     func(s string) string { return string(re.ReplaceAll([]byte(s), []byte("_"))) },
		"datef":  func(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) },
		"datef2": func(t time.Time) string { return t.UTC().Format(time.DateOnly) },
		"base":   func(s string) string { return path.Base(s) },
		"ext":    func(s string) string { return path.Ext(s) },
		"mime":   func(s string) string { return mime.TypeByExtension(path.Ext(s)) },
	})
	tmpl = template.Must(tmpl.ParseFS(tmplFS, "templates/*"))
}
