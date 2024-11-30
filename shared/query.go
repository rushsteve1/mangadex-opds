package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
	"time"
)

// APIUrl is the default MangaDex API URL
var APIUrl = url.URL{
	Scheme: "https",
	Host:   "api.mangadex.org",
}

// DevUrl is the MangaDex Dev API URL used in place of [APIUrl]
var DevUrl = url.URL{
	Scheme: "https",
	Host:   "api.mangadex.dev",
}

// UploadsURL is the MangaDex Uploads URL used when the MDUploads option is true
var UploadsURL = url.URL{
	Scheme: "https",
	Host:   "uploads.mangadex.org",
}

// UserAgent constructs the `User-Agent` header from the build information.
func UserAgent() string {
	info, ok := debug.ReadBuildInfo()

	// I have no idea under what circumstances this is possible but
	// defensive programming is the way to go
	if !ok {
		slog.Error("could not read build info")
		panic("could not read build info")
	}

	return fmt.Sprintf("%s/%s", path.Base(info.Main.Path), info.Main.Version)
}

// QueryAPILimiter limits the rate that QueryAPI is called
var QueryAPILimiter = rate.NewLimiter(rate.Every(time.Second)/5, 5)

// QueryAPI is used to fetch data from the MangaDex API.
func QueryAPI[T any](
	ctx context.Context,
	queryPath string,
	queryParams url.Values,
	limiter *rate.Limiter,
) (out T, err error) {
	var queryUrl url.URL
	if GlobalOptions.DevApi {
		queryUrl = DevUrl
	} else {
		queryUrl = APIUrl
	}

	queryUrl.Path = queryPath
	queryUrl.RawQuery = queryParams.Encode()

	var entry []byte

	for i := 0; i < 3 && err != nil; i++ {
		err = QueryAPILimiter.Wait(ctx)
		if err != nil {
			continue
		}
		if limiter != nil {
			err = limiter.Wait(ctx)
			if err != nil {
				continue
			}
		}

		slog.InfoContext(ctx, "querying API", "url", queryUrl.String())

		req, err2 := makeRequest(ctx, &queryUrl)
		if err2 != nil {
			err2 = err
			continue
		}
		req.Header.Set("Accept", "application/json")

		res, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			err2 = err
			continue
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("upstream error: %s", res.Status)
			time.After(1 * time.Minute)
			continue
		}

		entry, err = io.ReadAll(res.Body)
		if err == nil {
			break
		}
	}
	if err != nil {
		return out, err
	}
	err = json.NewDecoder(bytes.NewReader(entry)).Decode(&out)
	return out, err
}

// QueryImageLimiter limits the rate that QueryImage is called
var QueryImageLimiter = rate.NewLimiter(rate.Every(time.Second)/5, 5)

// QueryImage is used to fetch an image from the given URL.
// Only PNG and JPG images are supported, for compatibility with downstream CBZ and EPUB formats.
func QueryImage(
	ctx context.Context,
	imgUrl *url.URL,
	w io.Writer,
	limiter *rate.Limiter,
) (err error) {
	// In some tests we do not actually want to download the files
	if GlobalOptions.NoDownload {
		slog.Warn("no-download option enabled", "url", imgUrl.String())
		return nil
	}

	for i := 0; i < 3 && err != nil; i++ {
		err = QueryImageLimiter.Wait(ctx)
		if err != nil {
			continue
		}
		if limiter != nil {
			err = limiter.Wait(ctx)
			if err != nil {
				continue
			}
		}
		slog.InfoContext(ctx, "querying image", "url", imgUrl.String())

		req, err := makeRequest(ctx, imgUrl)
		if err != nil {
			return err
		}

		req.Header.Add("Accept", mime.TypeByExtension(".png"))
		req.Header.Add("Accept", mime.TypeByExtension(".jpg"))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("upstream error: %s", res.Status)
		}

		_, err = io.Copy(w, res.Body)
		if err == nil {
			break
		}
	}

	slog.DebugContext(ctx, "finished image download", "url", imgUrl)

	return err
}

func makeRequest(ctx context.Context, queryUrl *url.URL) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent())

	return req, nil
}

// TODO make this more general so it can be used for Chapters
func WithDefaultParams(queryParams url.Values) url.Values {
	if queryParams == nil {
		queryParams = url.Values{}
	}

	// Use reference expansion
	// https://api.mangadex.org/docs/01-concepts/reference-expansion/
	// TODO optimize these
	defaultParams := url.Values{
		"includes[]": []string{"author", "artist", "cover_art"},
	}

	for k, v := range defaultParams {
		queryParams[k] = v
	}

	return queryParams
}
