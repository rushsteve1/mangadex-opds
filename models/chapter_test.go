package models

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// Girl's Last Tour chapter 43 uploaded by rozen
const ChapterID = "9a612118-1441-431a-979d-85958fb20cf2"

func Test_FetchChapter(t *testing.T) {
	ctx := context.Background()

	c, err := FetchChapter(ctx, uuid.MustParse(ChapterID), nil)
	if err != nil {
		t.Fatal(err)
	}

	if c.Attributes.Title == "" {
		t.Fatal("no title")
	}

	t.Run("cast manga", func(t *testing.T) {
		m := c.Manga()
		if m == nil {
			t.Fatal("manga did not cast")
		}
	})
}

func Test_FetchImageURLs(t *testing.T) {
	ctx := context.Background()
	c := Chapter{ID: uuid.MustParse(ChapterID)}

	imgUrls, err := c.FetchImageURLs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(imgUrls) == 0 {
		t.Fatal("no image urls")
	}
}
