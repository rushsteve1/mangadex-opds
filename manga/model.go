package manga

import (
	"cmp"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/shared"
)

type Manga struct {
	ID            uuid.UUID             `json:"id"`
	Attributes    MangaAttributes       `json:"attributes"`
	Relationships []shared.Relationship `json:"relationships"`
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

func (m Manga) TrTitle() string {
	return shared.Tr(m.Attributes.Title)
}

func (m Manga) TrDesc() string {
	return shared.Tr(m.Attributes.Description)
}

type RelData struct {
	Authors  []Author
	CoverURL *url.URL
}

type Author struct {
	Name string
	URI  *url.URL
}

func ParseAuthor(m map[string]any) Author {
	u := cmp.Or(m["website"], m["twitter"], m["youtube"], m["tumblr"], m["pixiv"])
	uri, _ := url.Parse(u.(string))
	return Author{
		Name: m["name"].(string),
		URI:  uri,
	}
}

func (m Manga) GetRelData() (rd RelData) {
	for _, rel := range m.Relationships {
		switch rel.Type {
		case "author":
			fallthrough
		case "artist":
			rd.Authors = append(rd.Authors, ParseAuthor(rel.Attributes))
		case "cover_art":
			u := shared.UploadsURL
			u.Path, _ = url.JoinPath("covers", m.ID.String(), rel.Attributes["fileName"].(string))
			rd.CoverURL = &u
		}
	}
	return rd
}
