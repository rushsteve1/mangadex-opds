package models

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/rushsteve1/mangadex-opds/shared"

	"golang.org/x/sync/errgroup"
)

// Relationship represents a relation between one MangaDex data entity and another.
type Relationship struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes"`
}

func (r Relationship) URL() url.URL {
	relUrl := shared.APIUrl
	relUrl.Path, _ = url.JoinPath(r.Type, r.ID)
	return relUrl
}

// CastRelationship is a special helper function that uses JSON marshaling to
// transform a [Relationship]'s attributes into another type.
//
// This function should be used with care.
func CastRelationship[T any](rel *Relationship) (out T, err error) {
	r, w := io.Pipe()

	eg := errgroup.Group{}

	eg.Go(func() error {
		err = json.NewEncoder(w).Encode(rel.Attributes)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		return nil
	})

	err = json.NewDecoder(r).Decode(&out)
	if err != nil {
		return out, err
	}

	err = eg.Wait()
	return out, err
}
