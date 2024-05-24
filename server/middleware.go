package server

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rushsteve1/mangadex-opds/shared"
)

// SlogMiddleware logs out information about incoming HTTP requests.
func SlogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(
			r.Context(),
			"http request",
			"method", r.Method,
			"url", r.URL.String(),
			"remote", r.RemoteAddr,
		)

		if r.URL.Host != "" && r.URL.Host != shared.GlobalOptions.Host.Host {
			slog.WarnContext(
				r.Context(),
				"request from non-matching host",
				"actual",
				r.URL.Host,
				"expected",
				shared.GlobalOptions.Host.Host,
			)
		}

		next(w, r)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipMiddleware compresses the response body using GZIP if the client supports it.
// This can be used to save a lot of bandwidth at the cost of slightly more CPU usage.
func GzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next(gzr, r)
	}
}
