package manga

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// Bocchi's Guide to MangaDex
const MangaID = "d1c0d3f9-f359-467c-8474-0b2ea8e06f3d"

func Test_Fetch(t *testing.T) {
	ctx := context.Background()
	m, err := Fetch(ctx, uuid.MustParse(MangaID), nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Attributes.Title) == 0 {
		t.Fatal("no titles")
	}
}

func Test_Feed(t *testing.T) {
	ctx := context.Background()
	m := Manga{ID: uuid.MustParse(MangaID)}

	chapters, err := m.Feed(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(chapters) == 0 {
		t.Fatal("no chapters")
	}
}
