package manga

import (
	"time"

	"github.com/rushsteve1/mangadex-opds/shared"
)

type Manga struct {
	ID         shared.UUID     `json:"id"`
	Attributes MangaAttributes `json:"attributes"`
	// Relationships []shared.Relationship `json:"relationships"`
}

type MangaAttributes struct {
	Title            map[string]string `json:"title"`
	Description      map[string]string `json:"description"`
	OriginalLanguage string            `json:"originalLanguage"`
	Status           string            `json:"status"`
	Demographic      string            `json:"publicationDemographic"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
}
