package server

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"strconv"

	"github.com/rushsteve1/mangadex-opds/chapter"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/manga/search", todo)
	mux.HandleFunc("/manga/{id}", todo)

	mux.HandleFunc("/chapter/{id}", todo)

	eh := AcceptsMiddleware(epubHandler, mime.TypeByExtension(".epub"))
	mux.HandleFunc("/chapter/{id}/epub", eh)

	// TODO I don't think this mimetype is quite right, add more
	ch := AcceptsMiddleware(cbzHandler, mime.TypeByExtension(".cbz"))
	mux.HandleFunc("/chapter/{id}/cbz", ch)

	outerMux := http.NewServeMux()
	outerMux.HandleFunc("/", SlogMiddleware(mux.ServeHTTP))

	return outerMux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("mangadex-opds"))
}

func todo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func chapterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("page") {
		imageHandler(w, r)
		return
	}

	todo(w, r)
}

// Implement support for OPDS-PSE 1.0
// https://github.com/anansi-project/opds-pse/blob/master/v1.0.md
func imageHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// By checking this now we avoid the trips to the API
	if page < 0 {
		http.Error(w, "page must be > 0", http.StatusBadRequest)
	}

	// TODO we probably want to be able to differentiate upstream errors from network errors
	c, err := chapter.Fetch(r.Context(), r.PathValue("id"), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imgURLs, err := c.FetchImageURLs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pages are zero-indexed
	if page >= len(imgURLs) {
		resp := fmt.Sprintf("page %d is out of bounds, max %d", page, len(imgURLs)-1)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}

	imgURL := imgURLs[page]
	w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(imgURL.Path)))

	err = shared.QueryImage(r.Context(), imgURL, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func epubHandler(w http.ResponseWriter, r *http.Request) {
	c, err := chapter.Fetch(r.Context(), r.PathValue("id"), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".epub"))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.epub"`, c.FullTitle()))

	// TODO etags and caching?

	err = c.WriteEpub(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func cbzHandler(w http.ResponseWriter, r *http.Request) {
	c, err := chapter.Fetch(r.Context(), r.PathValue("id"), r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We're cheating a bit, the EPUBs that are created are also valid CBZs!
	// So as long as we tell the client that everything should work!
	w.Header().Set("Content-Type", mime.TypeByExtension(".cbz")) // TODO wrong mimetype?
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.cbz"`, c.FullTitle()))

	// TODO etags and caching?

	err = c.WriteEpub(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
