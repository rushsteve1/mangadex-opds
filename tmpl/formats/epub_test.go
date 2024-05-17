package formats

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/models"
)

// Girl's Last Tour chapter 43 uploaded by rozen
const ChapterID = "9a612118-1441-431a-979d-85958fb20cf2"

func Test_WriteEpub(t *testing.T) {
	ctx := context.Background()
	c := models.Chapter{ID: uuid.MustParse(ChapterID)}

	file, err := os.CreateTemp("", "*.epub")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	t.Log("temp file created at: " + file.Name())

	err = WriteEpub(ctx, &c, file)
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
