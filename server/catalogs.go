package server

import (
	"fmt"
	"net/http"

	"github.com/rushsteve1/mangadex-opds/manga"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := rootTemplate(w)
	if err != nil {
		die(w, r, err)
	}
}

func catalogSearchHandler(term string, order string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		params.Add(fmt.Sprintf("order[%s]", term), order)

		m, err := manga.Search(r.Context(), params)
		if err != nil {
			die(w, r, err)
			return
		}

		err = manga.MangaListFeed(w, m, r.URL.Path)
		if err != nil {
			die(w, r, err)
		}
	}
}

var newCatalogHandler = catalogSearchHandler("createdAt", "desc")
var popularCatalogHandler = catalogSearchHandler("followedCount", "desc")
var updatedCatalogHandler = catalogSearchHandler("latestUploadedChapter", "desc")
