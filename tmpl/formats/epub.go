package formats

import (
	"archive/zip"
	"context"
	"embed"
	"fmt"
	"github.com/rushsteve1/mangadex-opds/shared"
	"github.com/rushsteve1/mangadex-opds/tmpl"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"path"

	"github.com/rushsteve1/mangadex-opds/models"
)

//go:embed epub
var epubFiles embed.FS

// WriteEpub will write an EPUB file for the current [Chapter] to the given [io.Writer].
func WriteEpub(ctx context.Context, c *models.Chapter, w io.Writer) (err error) {
	z := zip.NewWriter(w)

	err = z.SetComment(c.FullTitle())
	if err != nil {
		return err
	}

	err = z.AddFS(epubFiles)
	if err != nil {
		return err
	}

	imgUrls, err := c.FetchImageURLs(ctx)
	if err != nil {
		return err
	}

	w, err = z.Create("content.opf")
	if err != nil {
		return err
	}

	err = tmpl.ContentOPF(c, w)
	if err != nil {
		return err
	}

	w, err = z.Create("toc.ncx")
	if err != nil {
		return err
	}

	err = tmpl.TocNCX(c, w)
	if err != nil {
		return err
	}

	imgChan := make(chan tmpl.ChapterImage)
	doneChan := make(chan error)

	// Fetch and add the image files in parallel
	go func() {
		eg, ctx := errgroup.WithContext(ctx)
		eg.SetLimit(3)

		for i, img := range imgUrls {
			eg.Go(func() error {
				imgName := path.Base(img.String())
				chImg := tmpl.ChapterImage{
					Name:  imgName,
					Index: i,
				}

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
		w, err = z.Create(fmt.Sprintf("%s.xhtml", img.Name))
		if err != nil {
			return err
		}

		err = tmpl.EpubXHTML(&img, w)
		if err != nil {
			return err
		}

		// Images will not be compressed, just stored
		// This saves a lot of time and performance at the cost of bigger files
		// But considering MangaDex is fine with hosting those I assume they're already optimized
		w, err = z.CreateHeader(&zip.FileHeader{
			Name:   img.Name,
			Method: zip.Store,
		})
		if err != nil {
			return err
		}

		_, err = io.Copy(w, &img.Data)
		if err != nil {
			return err
		}
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
