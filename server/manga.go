package server

import (
	"mime"
	"net/http"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/manga"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".xml"))

	// Return the opensearch XML document
	if len(r.URL.Query()) == 0 {
		data, err := tmplFS.ReadFile("templates/opensearch.xml")
		if err != nil {
			die(w, r, err)
			return
		}

		w.Write(data)
		return
	}

	// Otherwise return the OPDS XML list for the search
	resp, err := manga.Search(r.Context(), r.URL.Query())
	if err != nil {
		die(w, r, err)
		return
	}

	err = manga.MangaListFeed(w, resp, r.URL.Path)
	if err != nil {
		die(w, r, err)
	}
}

func mangaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
	}

	m, err := manga.Fetch(r.Context(), id, r.URL.Query())
	if err != nil {
		die(w, r, err)
		return
	}

	manga.MangaChapterFeed(r.Context(), w, m, r.URL.Query())
}
