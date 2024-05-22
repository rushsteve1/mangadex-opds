package formats

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"log/slog"
	"path"

	"github.com/rushsteve1/mangadex-opds/models"
	"github.com/rushsteve1/mangadex-opds/shared"
	"github.com/rushsteve1/mangadex-opds/tmpl"
	
	"golang.org/x/sync/errgroup"
)

type chapterImage struct {
	Index int
	Name  string
	Data  bytes.Buffer
}

func WriteCBZ(ctx context.Context, c *models.Chapter, w io.Writer) (err error) {
	z := zip.NewWriter(w)

	err = z.SetComment(c.FullTitle())
	if err != nil {
		return err
	}

	imgUrls, err := c.FetchImageURLs(ctx)
	if err != nil {
		return err
	}

	w, err = z.Create("ComicInfo.xml")
	if err != nil {
		return err
	}

	tmpl.ComicInfoXML(c, w)

	imgChan := make(chan chapterImage)
	doneChan := make(chan error)

	// Fetch and add the image files in parallel
	go func() {
		eg, ctx := errgroup.WithContext(ctx)
		eg.SetLimit(3)

		for _, img := range imgUrls {
			eg.Go(func() error {
				imgName := path.Base(img.String())
				chImg := chapterImage{Name: imgName}

				err := shared.QueryImage(ctx, img, &chImg.Data)
				if err != nil {
					return err
				}

				imgChan <- chImg

				return nil
			})
		}

		// Wait for all downloads to finish
		err = eg.Wait()
		close(imgChan)
		doneChan <- err

		slog.InfoContext(ctx, "done downloading images", "count", len(imgUrls))
	}()

	for img := range imgChan {
		// Images will not be compressed, just stored
		// This saves a lot of time and performance at the cost of bigger files
		// But considering MangaDex is fine with hosting those I assume they're already optimized
		w, err = z.CreateHeader(&zip.FileHeader{Name: img.Name, Method: zip.Store})
		if err != nil {
			return err
		}

		io.Copy(w, &img.Data)
	}

	err = <-doneChan
	if err != nil {
		return err
	}

	err = z.Close()
	if err != nil {
		return err
	}

	return nil
}
