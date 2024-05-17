package tmpl

import (
	"io"

	"github.com/rushsteve1/mangadex-opds/models"
)

func ComicInfoXML(c *models.Chapter, w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "comicinfo.tmpl.xml", c)
}
