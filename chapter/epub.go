package chapter

import (
	"archive/zip"
	"context"
	"io"
)

// WriteEpub will write an EPUB file for the current [Chapter] to the given [io.Writer].
func (c Chapter) WriteEpub(ctx context.Context, w io.Writer) (err error) {
	// A CBZ is almost a valid EPUB so we can use that as the base
	c.WriteCBZ(ctx, w)

	// err = c.writeEpubMetadata(z)

	return err
}

// TODO implement this function
// - Epub metadata
// - Epub structure
func (c Chapter) writeEpubMetadata(z *zip.Writer) (err error) {
	w, err := z.Create("META-INF/container.xml")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("TODO"))

	return err
}
