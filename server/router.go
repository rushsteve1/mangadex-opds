package server

import (
	"log/slog"
	"mime"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	ih := AcceptMiddleware(indexHandler, mime.TypeByExtension(".html"))
	mux.HandleFunc("/", ih)

	mux.HandleFunc("/root", AcceptXML(rootHandler))

	mux.HandleFunc("/search", AcceptXML(searchHandler))

	mux.HandleFunc("/catalogs/new", AcceptXML(newCatalogHandler))
	mux.HandleFunc("/catalogs/popular", AcceptXML(popularCatalogHandler))
	mux.HandleFunc("/catalogs/updated", AcceptXML(updatedCatalogHandler))

	mux.HandleFunc("/manga/{id}", AcceptXML(mangaHandler))

	eh := AcceptMiddleware(epubHandler, mime.TypeByExtension(".epub"))
	mux.HandleFunc("/chapter/{id}/epub", eh)

	// TODO I don't think this mimetype is quite right, add more
	ch := AcceptMiddleware(cbzHandler, mime.TypeByExtension(".cbz"))
	mux.HandleFunc("/chapter/{id}/cbz", ch)

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
