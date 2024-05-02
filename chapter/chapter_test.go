package chapter

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

// Girl's Last Tour chapter 43 uploaded by rozen
const ChapterID = "9a612118-1441-431a-979d-85958fb20cf2"

func Test_Fetch(t *testing.T) {
	ctx := context.Background()

	c, err := Fetch(ctx, uuid.MustParse(ChapterID), nil)
	if err != nil {
		t.Fatal(err)
	}

	if c.Attributes.Title == "" {
		t.Fatal("no title")
	}
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

func Test_WriteEpub(t *testing.T) {
	ctx := context.Background()
	c := Chapter{ID: uuid.MustParse(ChapterID)}

	file, err := os.CreateTemp("", "*.epub")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	t.Log("temp file created at: " + file.Name())

	err = c.WriteEpub(ctx, file)
	if err != nil {
		t.Fatal(err)
	}

	stats, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if stats.Size() == 0 {
		t.Fatal("size zero")
	}
}
