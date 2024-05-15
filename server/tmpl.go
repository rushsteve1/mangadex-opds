package server

import (
	"embed"
	"io"
	"text/template"
	"time"

	"github.com/rushsteve1/mangadex-opds/shared"
)

//go:embed templates
var tmplFS embed.FS
var tmpl = template.Must(template.ParseFS(tmplFS, "templates/*"))

//go:embed favicon.ico
var favicon []byte

//go:embed robots.txt
var robotstxt []byte

type indexData struct {
	Host    string
	Version string
}

func indexTemplate(w io.Writer) error {
	data := indexData{
		Host:    shared.GlobalOptions.Host.String(),
		Version: shared.Version,
	}

	return tmpl.ExecuteTemplate(w, "index.tmpl.html", data)
}

type rootData struct {
	UpdatedAt string
	Host      string
	Version   string
}

func rootTemplate(w io.Writer) error {
	data := rootData{
		UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Host:      shared.GlobalOptions.Host.String(),
		Version:   shared.Version,
	}

	return tmpl.ExecuteTemplate(w, "root.tmpl.xml", data)
}
