package server

import (
	"fmt"
	"net/http"

	"github.com/rushsteve1/mangadex-opds/manga"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := rootTemplate(w)
	if err != nil {
		httpError(w, r, err)
	}
}

func catalogSearchHandler(id string, title string, term string, order string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

var newCatalogHandler = catalogSearchHandler(
	"new",
	"New Manga",
	"createdAt",
	"desc",
)
var popularCatalogHandler = catalogSearchHandler(
	"popular",
	"Popular Manga",
	"followedCount",
	"desc",
)
var updatedCatalogHandler = catalogSearchHandler(
	"updated",
	"Recently Updated Manga",
	"latestUploadedChapter",
	"desc",
)
