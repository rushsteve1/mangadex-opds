package manga

import (
	"embed"
	"io"
	"text/template"
	"time"

	"github.com/google/uuid"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

func init() {
	tmpl = template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"datef": datef,
	})
	tmpl = template.Must(tmpl.ParseFS(tmplFS, "templates/*"))
}

type listData struct {
	ID        uuid.UUID
	SelfPath  string
	MangaList []Manga
}

func MangaListFeed(w io.Writer, mangaList []Manga, selfPath string) error {
	data := listData{
		ID:        uuid.New(),
		MangaList: mangaList,
		SelfPath:  selfPath,
	}

	return tmpl.ExecuteTemplate(w, "list.tmpl.xml", data)
}

func MangaChapterFeed(w io.Writer, m Manga) {}

func datef(t time.Time) string {
	return t.Format(time.RFC3339)
}
