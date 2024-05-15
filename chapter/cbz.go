package chapter

import (
	"archive/zip"
	"context"
	"io"
	"path"

	"github.com/rushsteve1/mangadex-opds/shared"
	"golang.org/x/sync/errgroup"
)

func (c Chapter) WriteCBZ(ctx context.Context, w io.Writer) (err error) {
	z := zip.NewWriter(w)

	err = z.SetComment(c.FullTitle())

	imgUrls, err := c.FetchImageURLs(ctx)
	imgNames := make([]string, len(imgUrls))

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(5)

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

	err = eg.Wait()
	if err != nil {
		return err
	}

	err = z.Close()
	if err != nil {
		return err
	}

	return nil
}
