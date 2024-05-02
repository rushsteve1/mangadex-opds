package shared

import (
	"net/url"
)

type Relationship struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes"`
}

func (r Relationship) URL() url.URL {
	relUrl := APIUrl
	relUrl.Path, _ = url.JoinPath(r.Type, r.ID)
	return relUrl
}
