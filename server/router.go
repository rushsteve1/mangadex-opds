package server

import (
	_ "embed"
	"expvar"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/rushsteve1/mangadex-opds/shared"
	"github.com/rushsteve1/mangadex-opds/tmpl"
)

//go:embed favicon.ico
var favicon []byte

//go:embed robots.txt
var robotstxt []byte

func init() {
	// Setup expvars
	expvar.Publish("buildinfo", (expvar.Func)(func() any {
		info, _ := debug.ReadBuildInfo()
		return info
	}))

	expvar.Publish("config", (expvar.Func)(func() any {
		return shared.GlobalOptions
	}))
}

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write(robotstxt)
	})

	mux.HandleFunc("/catalog", rootHandler)
	mux.HandleFunc("/catalog/new", newCatalogHandler)
	mux.HandleFunc("/catalog/popular", popularCatalogHandler)
	mux.HandleFunc("/catalog/updated", updatedCatalogHandler)

	mux.HandleFunc("/search", searchHandler)

	mux.HandleFunc("/manga/{id}", mangaHandler)

	mux.HandleFunc("/chapter/{id}", chapterHandler)
	mux.HandleFunc("/chapter/{id}/cbz", cbzHandler)
	mux.HandleFunc("/chapter/{id}/epub", epubHandler)

	mux.HandleFunc("/download", downloadHandler)

	// Panels hits this endpoint even when giving it a MD cover URL
	mux.HandleFunc("/covers/", coversHandler)

	// Enable expvars endpoint
	if shared.GlobalOptions.ExpVars {
		mux.Handle("/debug/vars", expvar.Handler())
	}

	outerMux := http.NewServeMux()
	innerMux := SlogMiddleware(mux.ServeHTTP)

	// This is on by default but can be turned off
	if shared.GlobalOptions.GzipResponses {
		innerMux = GzipMiddleware(innerMux)
	}

	outerMux.HandleFunc("/", innerMux)

	return outerMux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// This is the 404 handler
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		slog.WarnContext(r.Context(), "404 response", "url", r.URL.String())
		return
	}

	err := tmpl.IndexTemplate(w)
	if err != nil {
		httpError(w, r, err)
	}
}

func httpError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(
		r.Context(),
		"internal server error",
		"error",
		err.Error(),
		"path",
		r.URL.Path,
	)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
