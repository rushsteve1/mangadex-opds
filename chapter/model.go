package chapter

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chapter struct {
	ID         uuid.UUID       `json:"id"`
	Attributes ChapterAttributes `json:"attributes"`
	// Relationships []shared.Relationship `json:"relationships"`
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

func (c Chapter) FullTitle() string {
	builder := strings.Builder{}

	// TODO manga title

	if c.Attributes.Title != "" {
		builder.WriteString(c.Attributes.Title)

	}

	if c.Attributes.Volume != "" {
		builder.WriteString(" - vol ")
		builder.WriteString(c.Attributes.Volume)
	}

	if c.Attributes.Chapter != "" {
		builder.WriteString(" - ch ")
		builder.WriteString(c.Attributes.Chapter)
	}

	return builder.String()
}
