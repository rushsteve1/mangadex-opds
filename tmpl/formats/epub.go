package formats

import (
	"archive/zip"
	"context"
	"io"

	"github.com/rushsteve1/mangadex-opds/models"
)

// WriteEpub will write an EPUB file for the current [Chapter] to the given [io.Writer].
func WriteEpub(ctx context.Context, c *models.Chapter, w io.Writer) (err error) {
	// A CBZ is almost a valid EPUB so we can use that as the base
	WriteCBZ(ctx, c, w)

	// err = c.writeEpubMetadata(z)

	return err
}

// TODO implement this function
// - Epub metadata
// - Epub structure
func writeEpubMetadata(c *models.Chapter, z *zip.Writer) (err error) {
	w, err := z.Create("META-INF/container.xml")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("TODO"))

	return err
}
