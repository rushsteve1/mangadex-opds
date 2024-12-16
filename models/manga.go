package models

import (
	"cmp"
	"net/url"
	"time"

	"github.com/rushsteve1/mangadex-opds/shared"

	"github.com/google/uuid"
)

// Manga is a single manga/comic series on MangaDex.
// Use [FetchManga] to get one from the API.
type Manga struct {
	ID            uuid.UUID       `json:"id"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
	relData       *RelData
}

// MangaAttributes are all the field attributes for a manga.
type MangaAttributes struct {
	Title            map[string]string   `json:"title"`
	AltTitles        []map[string]string `json:"altTitles"`
	Description      map[string]string   `json:"description"`
	OriginalLanguage string              `json:"originalLanguage"`
	Status           string              `json:"status"`
	Demographic      string              `json:"publicationDemographic"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

func (m *Manga) URL() string {
	u := shared.GlobalOptions.Host
	u.Path, _ = url.JoinPath("manga", m.ID.String())
	return u.String()
}

func (m *Manga) mergeTitles() {
	for _, at := range m.Attributes.AltTitles {
		for k, v := range at {
			_, found := m.Attributes.Title[k]
			if !found {
				m.Attributes.Title[k] = v
			}
		}
	}
}

// TrTitle returns the [Manga]'s title translated using [shared.Tr].
func (m Manga) TrTitle() string {
	return shared.Tr(m.Attributes.Title)
}

// TrDesc retuns the [Manga]'s description translated using [shared.Tr]
func (m Manga) TrDesc() string {
	return shared.Tr(m.Attributes.Description)
}

// RelData is the parsed [Relationship] data for a [Manga].
type RelData struct {
	Authors  []Author
	CoverURL string
}

// Author is a single author/artist for a [Manga],
// of which there may be several.
type Author struct {
	Name string
	URI  string
}

func parseAuthor(m map[string]any) (a Author) {
	// MangaDex exposes a variety of options for this
	ustr := cmp.Or(m["website"], m["twitter"], m["youtube"], m["tumblr"], m["pixiv"])
	s, ok := ustr.(string)
	if ok {
		u, _ := url.Parse(s)
		a.URI = u.String()
	}

	s, ok = m["name"].(string)
	if ok {
		a.Name = s
	} else {
		a.Name = "Unknown Author"
	}

	return a
}

// RelData returns the [RelData] for this manga.
// This value is cached on the [Manga] struct.
func (m *Manga) RelData() (rd *RelData) {
	if m.relData != nil {
		return m.relData
	}

	rd = &RelData{}

	for _, rel := range m.Relationships {
		switch rel.Type {
		case "author":
			fallthrough
		case "artist":
			// TODO dedup same author/artist
			rd.Authors = append(rd.Authors, parseAuthor(rel.Attributes))
		case "cover_art":
			u := shared.UploadsURL
			fn, ok := rel.Attributes["fileName"].(string)
			if ok {
				u.Path, _ = url.JoinPath("covers", m.ID.String(), fn)
				rd.CoverURL = u.String()
			}
		}
	}

	m.relData = rd
	return m.relData
}
