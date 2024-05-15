package server

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/rushsteve1/mangadex-opds/manga"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

	err := rootTemplate(w)
	if err != nil {
		httpError(w, r, err)
	}
}

func makeCatalogHandler(id string, title string, term string, order string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

		params := r.URL.Query()
		params.Add(fmt.Sprintf("order[%s]", term), order)

		m, err := manga.Search(r.Context(), params)
		if err != nil {
			httpError(w, r, err)
			return
		}

		err = manga.MangaListFeed(w, id, title, m, r.URL.Path)
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
