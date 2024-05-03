package server

import (
	"log/slog"
	"mime"
	"net/http"
	"strings"
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

		next(w, r)
	}
}

func HtmlXmlSplitterMiddleware(htmlNext http.HandlerFunc, xmlNext http.HandlerFunc) http.HandlerFunc {
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
