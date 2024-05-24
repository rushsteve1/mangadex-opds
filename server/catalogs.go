package server

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/rushsteve1/mangadex-opds/models"
	"github.com/rushsteve1/mangadex-opds/tmpl"
)

// rootHandler returns the static root catalog [tmpl.RootTemplate]
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

	err := tmpl.RootTemplate(w)
	if err != nil {
		httpError(w, r, err)
	}
}

// makeCatalogHandler creates HTTP catalog handlers for various manga searches
func makeCatalogHandler(id string, title string, term string, order string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

		params := r.URL.Query()
		params.Add(fmt.Sprintf("order[%s]", term), order)

		m, err := models.SearchManga(r.Context(), params)
		if err != nil {
			httpError(w, r, err)
			return
		}

		err = tmpl.MangaListFeed(w, id, title, m, r.URL.Path)
		if err != nil {
			httpError(w, r, err)
		}
	}
}

var newCatalogHandler = makeCatalogHandler(
	"new",
	"New Manga",
	"createdAt",
	"desc",
)
var popularCatalogHandler = makeCatalogHandler(
	"popular",
	"Popular Manga",
	"followedCount",
	"desc",
)
var updatedCatalogHandler = makeCatalogHandler(
	"updated",
	"Recently Updated Manga",
	"latestUploadedChapter",
	"desc",
)
