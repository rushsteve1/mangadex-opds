package tmpl

import (
	"io"
	"time"

	"github.com/rushsteve1/mangadex-opds/shared"
)

func OpenSearchXML(w io.Writer) error {
	f, err := tmplFS.Open("templates/opensearch.xml")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}

type indexData struct {
	Host    string
	Version string
}

func IndexTemplate(w io.Writer) error {
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

func RootTemplate(w io.Writer) error {
	data := rootData{
		UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Host:      shared.GlobalOptions.Host.String(),
		Version:   shared.Version,
	}

	return tmpl.ExecuteTemplate(w, "root.tmpl.xml", data)
}
