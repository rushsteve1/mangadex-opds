package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rushsteve1/mangadex-opds/shared"
)

func Test_PathsOK(t *testing.T) {
	shared.TestOptions()

	srv := httptest.NewServer(Router())

	testCases := []string{
		"/",
		"/catalog",
		"/search",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			fullUrl := fmt.Sprintf("%s%s", srv.URL, tc)

			t.Logf("full url: %s", fullUrl)

			resp, err := srv.Client().Get(fullUrl)
			shared.AssertEq(t, err, nil)
			shared.AssertNeq(t, resp, nil)
			shared.AssertEq(t, resp.StatusCode, http.StatusOK)
		})
	}
}

func Test_DownloadMDUrl(t *testing.T) {
	shared.TestOptions()
	shared.GlobalOptions.NoDownload = true

	testCases := []string{
		"https://mangadex.org/chapter/9a612118-1441-431a-979d-85958fb20cf2",
		"http://mangadex.org/chapter/9a612118-1441-431a-979d-85958fb20cf2/2",
		"http://mangadex.cc/chapter/9a612118-1441-431a-979d-85958fb20cf2/2",
		"http://www.mangadex.org/chapter/9a612118-1441-431a-979d-85958fb20cf2",
	}

	srv := httptest.NewServer(Router())

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			fullUrl := fmt.Sprintf("%s/download?url=%s", srv.URL, url.QueryEscape(tc))
			t.Logf("full url: %s", fullUrl)

			resp, err := srv.Client().Get(
				fullUrl,
			)
			shared.AssertEq(t, err, nil)
			shared.AssertNeq(t, resp, nil)
			defer resp.Body.Close()

			shared.AssertEq(t, resp.StatusCode, http.StatusOK)

			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			shared.AssertEq(t, err, nil)
			shared.AssertEq(t, len(body), 4878)
		})
	}
}
