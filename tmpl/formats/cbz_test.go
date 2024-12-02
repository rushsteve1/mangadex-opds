package formats

import (
	"bytes"
	"context"
	"testing"

	"github.com/rushsteve1/mangadex-opds/models"
	"github.com/rushsteve1/mangadex-opds/shared"

	"github.com/google/uuid"
)

// Girl's Last Tour chapter 43 uploaded by rozen
const ChapterID = "9a612118-1441-431a-979d-85958fb20cf2"

// We use size to check validity because we can't use hashes due to ModTime
// and it is extremely unlikely that an invalid zip would have exactly the right size
const CBZChapterSize = 1_173_566

func Test_WriteCBZ(t *testing.T) {
	shared.TestOptions()

	ctx := context.Background()

	c, err := models.FetchChapter(ctx, uuid.MustParse(ChapterID), nil)
	shared.AssertEq(t, err, nil)

	buf := bytes.Buffer{}

	err = WriteCBZ(ctx, &c, &buf)
	shared.AssertEq(t, err, nil)
	shared.AssertEq(t, buf.Len(), CBZChapterSize)
}
