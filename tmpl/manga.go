package tmpl

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/rushsteve1/mangadex-opds/models"
	"github.com/rushsteve1/mangadex-opds/shared"
)

type listData struct {
	ID        string
	Self      string
	MangaList []models.Manga
	Host      string
	UpdatedAt string
	Title     string
	Version   string
}

func MangaListFeed(
	w io.Writer,
	id string,
	title string,
	mangaList []models.Manga,
	selfPath string,
) error {
	self := shared.GlobalOptions.Host
	self.Path = selfPath
	data := listData{
		ID:        id,
		MangaList: mangaList,
		Self:      self.String(),
		Host:      shared.GlobalOptions.Host.String(),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Title:     title,
		Version:   shared.Version,
	}

	return tmpl.ExecuteTemplate(w, "list.tmpl.xml", data)
}

type chaptersData struct {
	Manga    *models.Manga
	Chapters []models.Chapter
	Host     string
	Version  string
}

func MangaChapterFeed(
	ctx context.Context,
	w io.Writer,
	m *models.Manga,
	queryParams url.Values,
) error {
	chapters, err := m.Feed(ctx, queryParams)
	if err != nil {
		return err
	}

	data := chaptersData{
		Manga:    m,
		Chapters: chapters,
		Host:     shared.GlobalOptions.Host.String(),
		Version:  shared.Version,
	}

	return tmpl.ExecuteTemplate(w, "chapters.tmpl.xml", data)
}
