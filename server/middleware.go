package server

import (
	"fmt"
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

func AcceptsMiddleware(next http.HandlerFunc, mtypes ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accepts := r.Header.Get("Accepts")
		if containsAny(accepts, mtypes...) {
			next(w, r)
		}

		resp := fmt.Sprintf("accepts header must contain one of %v", mtypes)
		http.Error(w, resp, http.StatusNotAcceptable)
	}
}

func containsAny(s string, ss ...string) bool {
	for _, str := range ss {
		if strings.Contains(s, str) {
			return true
		}
	}

	return false
}
