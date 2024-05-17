package models

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/shared"
)

// Girls' Last Tour
const MangaID = "5b93fa0f-0640-49b8-974e-954b9959929b"

func Test_FetchManga(t *testing.T) {
	shared.GlobalOptions = shared.TestOptions()

	testCases := []struct {
		Title string
		ID    string
	}{
		{
			Title: "Girls' Last Tour",
			ID:    MangaID,
		},
		{
			Title: "Oyasumi Punpun",
			ID:    "4301d363-ee02-43f4-ae24-4cbf29a74830",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Title, func(t *testing.T) {
			ctx := context.Background()
			m, err := FetchManga(ctx, uuid.MustParse(tc.ID), nil)
			if err != nil {
				t.Fatal(err)
			}

			if len(m.Attributes.Title) == 0 {
				t.Fatal("no titles")
			}

			t.Run("test TrTitle", func(t *testing.T) {
				title := m.TrTitle()
				if title != tc.Title {
					t.Fatalf("TrTitle did not match expected: '%s' actual '%s'", tc.Title, title)
				}
			})
		})
	}
}

func Test_Feed(t *testing.T) {
	shared.GlobalOptions = shared.TestOptions()

	shared.GlobalOptions = shared.TestOptions()

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
