package models

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/rushsteve1/mangadex-opds/shared"

	"github.com/google/uuid"
)

// FetchManga gets the manga series information the MangaDex API and returns the [Manga].
func FetchManga(ctx context.Context, id uuid.UUID, queryParams url.Values) (m Manga, err error) {
	slog.InfoContext(ctx, "fetching manga", "id", id)

	queryPath, err := url.JoinPath("manga", id.String())
	if err != nil {
		return m, err
	}

	queryParams = shared.WithDefaultParams(queryParams)

	data, err := shared.QueryAPI[Data[Manga]](ctx, queryPath, queryParams)
	if err != nil {
		return m, err
	}

	m = data.Data
	m.mergeTitles()
	m.RelData()

	return m, err
}

// SearchManga queries the MangaDex search endpoint and returns an array of [Manga].
func SearchManga(ctx context.Context, queryParams url.Values) (ms []Manga, err error) {
	queryParams = shared.WithDefaultParams(queryParams)

	data, err := shared.QueryAPI[Data[[]Manga]](ctx, "manga", queryParams)
	if err != nil {
		return nil, err
	}

	for i := range data.Data {
		data.Data[i].RelData()
	}

	return data.Data, err
}

// Feed returns the chapter feed for a series as an array of [Chapter].
// By default the it filters to the current language in [shared.GlobalOptions]
// and sorts the chapters in ascending order, filtering out empty chapters.
// This can be changed using the queryParams.
func (m *Manga) Feed(ctx context.Context, queryParams url.Values) (cs []Chapter, err error) {
	queryPath, err := url.JoinPath("manga", m.ID.String(), "feed")
	if err != nil {
		return nil, err
	}

	if queryParams == nil {
		queryParams = url.Values{}
	}

	queryParams.Add("order[chapter]", "asc")
	queryParams.Add("translatedLanguage[]", shared.GlobalOptions.Language)
	queryParams.Add("includeEmptyPages", "0")

	data, err := shared.QueryAPI[Data[[]Chapter]](ctx, queryPath, queryParams)

	for i := range data.Data {
		data.Data[i].manga = m
		data.Data[i].FullTitle()
	}

	return data.Data, err
}
