package server

import (
	"embed"
	"io"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/shared"
)

//go:embed templates
var tmplFS embed.FS
var tmpl = template.Must(template.ParseFS(tmplFS, "templates/*"))

type indexData struct {
	Host string
}

func indexTemplate(w io.Writer) error {
	data := indexData{
		Host: shared.GlobalOptions.Host.String(),
	}

	return tmpl.ExecuteTemplate(w, "index.tmpl.html", data)
}

type rootData struct {
	ID        uuid.UUID
	UpdatedAt string
	Host      string
}

func rootTemplate(w io.Writer) error {
	data := rootData{
		ID:        uuid.New(),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Host:      shared.GlobalOptions.Host.String(),
	}

	return tmpl.ExecuteTemplate(w, "root.tmpl.xml", data)
}
