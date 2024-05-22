package models

import (
	"context"
	"testing"

	"github.com/rushsteve1/mangadex-opds/shared"

	"github.com/google/uuid"
)

// Girl's Last Tour chapter 43 uploaded by rozen
const ChapterID = "9a612118-1441-431a-979d-85958fb20cf2"

func Test_FetchChapter(t *testing.T) {
	shared.TestOptions()

	ctx := context.Background()

	c, err := FetchChapter(ctx, uuid.MustParse(ChapterID), nil)
	shared.AssertEq(t, err, nil)
	shared.AssertNeq(t, c.Attributes.Title, "")

	t.Run("cast manga", func(t *testing.T) {
		m := c.Manga()
		shared.AssertNeq(t, m, nil)
	})
}

func Test_FetchImageURLs(t *testing.T) {
	shared.TestOptions()

	ctx := context.Background()
	c := Chapter{ID: uuid.MustParse(ChapterID)}

	imgUrls, err := c.FetchImageURLs(ctx)
	shared.AssertEq(t, err, nil)
	shared.AssertNeq(t, len(imgUrls), 0)
}
