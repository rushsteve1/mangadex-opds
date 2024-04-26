package chapter

import (
	"archive/zip"
	"context"
	"io"
	"path"

	"golang.org/x/sync/errgroup"

	"github.com/rushsteve1/mangadex-opds/shared"
)

// WriteEpub will write an EPUB file for the current [Chapter] to the given [io.Writer].
// As a bit of a trick the generated EPUBs are also valid CBZ files.
func (c Chapter) WriteEpub(ctx context.Context, w io.Writer) (err error) {
	z := zip.NewWriter(w)
	defer z.Close()

	err = z.SetComment(c.FullTitle())

	imgUrls, err := c.FetchImageURLs(ctx)
	imgNames := make([]string, len(imgUrls))

	eg, ctx := errgroup.WithContext(ctx)

	// Fetch and add the image files in parallel
	for i, img := range imgUrls {
		eg.Go(func() error {
			imgNames[i] = path.Base(img.String())

			// Images will not be compressed, just stored
			// This saves a lot of time and performance at the cost of bigger EPUB files
			// But considering MangaDex is fine with hosting those I assume they're already optimized
			w, err = z.CreateHeader(&zip.FileHeader{Name: imgNames[i], Method: zip.Store})
			if err != nil {
				return err
			}

			return shared.QueryImage(ctx, img, w)
		})
	}

	err = c.writeMetadata(z)

	return err
}

func (c Chapter) writeMetadata(z *zip.Writer) (err error) {
	w, err := z.Create("META-INF/container.xml")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("TODO"))

	return err
}
