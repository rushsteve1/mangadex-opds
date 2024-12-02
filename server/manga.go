package server

import (
	"mime"
	"net/http"

	"github.com/rushsteve1/mangadex-opds/models"
	"github.com/rushsteve1/mangadex-opds/tmpl"

	"github.com/google/uuid"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

	// Return the opensearch XML document
	if len(r.URL.Query()) == 0 {
		err := tmpl.OpenSearchXML(w)
		if err != nil {
			httpError(w, r, err)
			return
		}
		return
	}

	// Otherwise return the OPDS XML list for the search
	resp, err := models.SearchManga(r.Context(), r.URL.Query())
	if err != nil {
		httpError(w, r, err)
		return
	}

	err = tmpl.MangaListFeed(w, "search", "Search Manga", resp, r.URL.Path)
	if err != nil {
		httpError(w, r, err)
	}
}

func mangaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
	}

	m, err := models.FetchManga(r.Context(), id, r.URL.Query())
	if err != nil {
		httpError(w, r, err)
		return
	}

	err = tmpl.MangaChapterFeed(r.Context(), w, &m, r.URL.Query())
	if err != nil {
		httpError(w, r, err)
	}
}
