package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
)

var APIUrl = url.URL{
	Scheme: "https",
	Host:   "api.mangadex.org",
}

var DevUrl = url.URL{
	Scheme: "https",
	Host:   "api.mangadex.dev",
}

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

// QueryAPI is used to fetch data from the MangaDex API.
func QueryAPI[T any](ctx context.Context, queryPath string, queryParams url.Values) (out T, err error) {
	var queryUrl url.URL
	if GlobalOptions.DevApi {
		queryUrl = DevUrl
	} else {
		queryUrl = APIUrl
	}

	queryUrl.Path = queryPath
	queryUrl.RawQuery = queryParams.Encode()

	slog.InfoContext(ctx, "querying API", "url", queryUrl.String())

	req, err := makeRequest(ctx, &queryUrl)
	if err != nil {
		return out, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return out, fmt.Errorf("upstream error: %s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}

// QueryImage is used to fetch an image from the given URL.
// Only PNG and JPG images are supported, for compatibility.
func QueryImage(ctx context.Context, imgUrl *url.URL, w io.Writer) (err error) {
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
