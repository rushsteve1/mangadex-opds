package server

import (
	"log/slog"
	"net/http"
)

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

	mux.HandleFunc("/chapter/{id}/epub", epubHandler)
	mux.HandleFunc("/chapter/{id}/cbz", cbzHandler)

	outerMux := http.NewServeMux()
	outerMux.HandleFunc("/", SlogMiddleware(mux.ServeHTTP))

	return outerMux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// This is the 404 handler
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err := indexTemplate(w)
	if err != nil {
		die(w, r, err)
	}
}

func die(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "internal server error", "error", err.Error(), "path", r.URL.Path)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
