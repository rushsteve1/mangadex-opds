package tmpl

import (
	"bytes"
	"io"

	"github.com/rushsteve1/mangadex-opds/models"
)

type ChapterImage struct {
	Index int
	Name  string
	Data  bytes.Buffer
}

func ComicInfoXML(c *models.Chapter, w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "comicinfo.tmpl.xml", c)
}

func ContentOPF(c *models.Chapter, w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "content.tmpl.opf", c)
}

func TocNCX(c *models.Chapter, w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "toc.tmpl.ncx", c)
}

func EpubXHTML(c *ChapterImage, w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "epub.tmpl.xhtml", c)
}
