package models

import (
	"cmp"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/shared"
)

// Chapter is a single chapter on MangaDex.
// Use [FetchChapter] to get one from the API.
type Chapter struct {
	ID            uuid.UUID         `json:"id"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []Relationship    `json:"relationships"`
	fullTitle     string
	manga         *Manga
	imgUrls       []*url.URL
}

// ChapterAttributes are all the field attributes for a chapter.
type ChapterAttributes struct {
	Title              string    `json:"title"`
	Volume             string    `json:"volume"`
	Chapter            string    `json:"chapter"`
	Pages              int       `json:"pages"`
	TranslatedLanguage string    `json:"translatedLanguage"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (c *Chapter) URL() string {
	u := shared.GlobalOptions.Host
	u.Path, _ = url.JoinPath("chapter", c.ID.String())
	return u.String()
}

// FullTitle builds and returns the full title of the chapter including
// the manga name, volume number, chapter number, and chapter title.
// This value is cached on the [Chapter] struct.
func (c *Chapter) FullTitle() string {
	if len(c.fullTitle) > 0 {
		return c.fullTitle
	}

	builder := strings.Builder{}

	m := c.Manga()
	var title string
	if m == nil {
		title = "Unknown Manga"
	} else {
		title = m.TrTitle()
	}

	builder.WriteString(title)
	builder.WriteString(" - ")

	if c.Attributes.Volume != "" {
		builder.WriteString("[Vol. ")
		builder.WriteString(c.Attributes.Volume)
		builder.WriteString("]")
	}

	builder.WriteString(" Chapter ")
	builder.WriteString(cmp.Or(c.Attributes.Chapter, "Unknown"))

	if c.Attributes.Title != "" {
		builder.WriteString(" - ")
		builder.WriteString(c.Attributes.Title)
	}

	c.fullTitle = strings.TrimSpace(builder.String())
	return c.fullTitle
}

// Manga finds and casts the [Manga] that this [Chapter] belongs to.
// This value is cached on the [Chapter] struct.
func (c *Chapter) Manga() *Manga {
	if c.manga != nil {
		return c.manga
	}

	for _, rel := range c.Relationships {
		if rel.Type == "manga" {
			c.manga = &Manga{}
			a, err := CastRelationship[MangaAttributes](&rel)
			if err != nil {
				slog.Error("error casting to manga", "error", err)
				return nil
			}
			c.manga.ID = uuid.MustParse(rel.ID)
			c.manga.Attributes = a
			c.manga.RelData()

			return c.manga
		}
	}

	return nil
}

func (c *Chapter) ImgURLs() []*url.URL {
	return c.imgUrls
}
