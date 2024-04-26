package manga

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/rushsteve1/mangadex-opds/chapter"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func Fetch(ctx context.Context, id shared.UUID, queryParams url.Values) (m Manga, err error) {
	slog.InfoContext(ctx, "fetching manga", "id", id)

	queryPath, err := url.JoinPath("manga", id)
	if err != nil {
		return m, err
	}

	if queryParams == nil {
		queryParams = url.Values{}
	}

	// Use reference expansion
	// https://api.mangadex.org/docs/01-concepts/reference-expansion/
	// TODO optimize these
	defaultParams := url.Values{
		"includes[]": []string{"author", "artist", "cover_art"},
	}

	for k, v := range defaultParams {
		queryParams[k] = v
	}

	data, err := shared.QueryAPI[shared.Data[Manga]](ctx, queryPath, queryParams)

	return data.Data, err
}

func Search(ctx context.Context, queryParams url.Values) (ms []Manga, err error) {
	data, err := shared.QueryAPI[shared.Data[[]Manga]](ctx, "manga", queryParams)

	return data.Data, err
}

func (m Manga) Feed(ctx context.Context, queryParams url.Values) (cs []chapter.Chapter, err error) {
	queryPath, err := url.JoinPath("manga", m.ID, "feed")
	if err != nil {
		return nil, err
	}

	data, err := shared.QueryAPI[shared.Data[[]chapter.Chapter]](ctx, queryPath, queryParams)

	return data.Data, err
}
