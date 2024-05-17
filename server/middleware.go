package server

import (
	"compress/gzip"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strings"

	"github.com/rushsteve1/mangadex-opds/shared"
)

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

func HtmlXmlSplitterMiddleware(
	htmlNext http.HandlerFunc,
	xmlNext http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header.Get("Accepts")

		if strings.Contains(accepts, mime.TypeByExtension(".html")) {
			htmlNext(w, r)
			return
		}

		if strings.Contains(accepts, mime.TypeByExtension(".xml")) {
			xmlNext(w, r)
			return
		}

		http.Error(w, "invalid accepts header", http.StatusNotAcceptable)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

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
