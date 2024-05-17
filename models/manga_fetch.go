package models

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func FetchManga(ctx context.Context, id uuid.UUID, queryParams url.Values) (m Manga, err error) {
	slog.InfoContext(ctx, "fetching manga", "id", id)

	queryPath, err := url.JoinPath("manga", id.String())
	if err != nil {
		return m, err
	}

	queryParams = shared.WithDefaultParams(queryParams)

	data, err := shared.QueryAPI[Data[Manga]](ctx, queryPath, queryParams)

	m = data.Data
	m.MergeTitles()

	return m, err
}

func Search(ctx context.Context, queryParams url.Values) (ms []Manga, err error) {
	queryParams = shared.WithDefaultParams(queryParams)

	data, err := shared.QueryAPI[Data[[]Manga]](ctx, "manga", queryParams)

	return data.Data, err
}

func (m Manga) Feed(ctx context.Context, queryParams url.Values) (cs []Chapter, err error) {
	queryPath, err := url.JoinPath("manga", m.ID.String(), "feed")
	if err != nil {
		return nil, err
	}

	if queryParams == nil {
		queryParams = url.Values{}
	}

	queryParams.Set("order[chapter]", "asc")
	queryParams.Set("translatedLanguage[]", shared.GlobalOptions.Language)
	queryParams.Set("includeEmptyPages", "0")

	data, err := shared.QueryAPI[Data[[]Chapter]](ctx, queryPath, queryParams)

	return data.Data, err
}
