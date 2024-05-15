package chapter

import (
	"cmp"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/shared"
)

type Chapter struct {
	ID            uuid.UUID             `json:"id"`
	Attributes    ChapterAttributes     `json:"attributes"`
	Relationships []shared.Relationship `json:"relationships"`
}

type ChapterAttributes struct {
	Title              string    `json:"title"`
	Volume             string    `json:"volume"`
	Chapter            string    `json:"chapter"`
	Pages              int       `json:"pages"`
	TranslatedLanguage string    `json:"translatedLanguage"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (c Chapter) URL() string {
	u := shared.GlobalOptions.Host
	u.Path, _ = url.JoinPath("chapter", c.ID.String())
	return u.String()
}

func (c Chapter) FullTitle() string {
	builder := strings.Builder{}

	// TODO manga title

	if c.Attributes.Volume != "" {
		builder.WriteString("[Vol. ")
		builder.WriteString(c.Attributes.Volume)
		builder.WriteString("]")
	}

	builder.WriteString(" Chapter ")
	builder.WriteString(cmp.Or(c.Attributes.Chapter, "Unknown"))

	if c.Attributes.Title != "" {
		if c.Attributes.Chapter != "" || c.Attributes.Volume != "" {
			builder.WriteString(" - ")
		}
		builder.WriteString(c.Attributes.Title)
	}

	return strings.TrimSpace(builder.String())
}
