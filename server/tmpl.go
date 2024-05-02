package server

import (
	"embed"
	"io"
	"text/template"
	"time"

	"github.com/google/uuid"
)

//go:embed templates
var tmplFS embed.FS
var tmpl = template.Must(template.ParseFS(tmplFS, "templates/*"))

type rootData struct {
	ID        uuid.UUID
	UpdatedAt string
}

func rootTemplate(w io.Writer) error {
	data := rootData{
		ID:        uuid.New(),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return tmpl.ExecuteTemplate(w, "root.tmpl.xml", data)
}
