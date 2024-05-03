package manga

import (
	"context"
	"embed"
	"io"
	"log/slog"
	"mime"
	"net/url"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/chapter"
	"github.com/rushsteve1/mangadex-opds/shared"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

func init() {
	tmpl = template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"datef": func(t time.Time) string { return t.Format(time.RFC3339) },
		"mime":  func(s string) string { return mime.TypeByExtension(path.Ext(s)) },
	})
	tmpl = template.Must(tmpl.ParseFS(tmplFS, "templates/*"))
}

type listData struct {
	ID        uuid.UUID
	Self      string
	MangaList []Manga
	Host      string
}

func MangaListFeed(w io.Writer, mangaList []Manga, selfPath string) error {
	self := shared.GlobalOptions.Host
	self.Path = selfPath
	data := listData{
		ID:        uuid.New(),
		MangaList: mangaList,
		Self:      self.String(),
		Host:      shared.GlobalOptions.Host.String(),
	}

	return tmpl.ExecuteTemplate(w, "list.tmpl.xml", data)
}

type chaptersData struct {
	Manga    Manga
	Chapters []chapter.Chapter
	Host     string
}

func MangaChapterFeed(ctx context.Context, w io.Writer, m Manga, queryParams url.Values) error {
	chapters, err := m.Feed(ctx, queryParams)
	if err != nil {
		return err
	}

	if len(chapters) == 0 {
		panic("fuck")
	} else {
		slog.InfoContext(ctx, "got chapters", "count", len(chapters))
	}

	data := chaptersData{
		Manga:    m,
		Chapters: chapters,
		Host:     shared.GlobalOptions.Host.String(),
	}

	return tmpl.ExecuteTemplate(w, "chapters.tmpl.xml", data)
}
