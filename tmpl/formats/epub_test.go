package formats

import (
	"bytes"
	"context"
	"github.com/rushsteve1/mangadex-opds/shared"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/models"
)

const EPUBChapterSize = 1_188_189

func Test_WriteEpub(t *testing.T) {
	shared.TestOptions()

	ctx := context.Background()

	c, err := models.FetchChapter(ctx, uuid.MustParse(ChapterID), nil)
	shared.AssertEq(t, err, nil)

	buf := bytes.Buffer{}

	err = os.Chdir("../../")
	err = WriteEpub(ctx, &c, &buf)
	shared.AssertEq(t, err, nil)
	shared.AssertEq(t, buf.Len(), EPUBChapterSize)
}
