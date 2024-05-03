package server

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/google/uuid"
	"github.com/rushsteve1/mangadex-opds/chapter"
	"github.com/rushsteve1/mangadex-opds/shared"
)

func init() {
	// I think this will help get around an issue with the mime types
	mime.AddExtensionType(".cbz", "application/zip")
}

func chapterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("page") {
		imageHandler(w, r)
		return
	}

	u := r.URL
	u.Path, _ = url.JoinPath(r.URL.Path, "epub")

	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}

// Implement support for OPDS-PSE 1.0
// https://github.com/anansi-project/opds-pse/blob/master/v1.0.md
func imageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		die(w, r, err)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		die(w, r, err)
		return
	}

	// By checking this now we avoid the trips to the API
	if page < 0 {
		die(w, r, err)
		return
	}

	// TODO we probably want to be able to differentiate upstream errors from network errors
	c, err := chapter.Fetch(r.Context(), id, r.URL.Query())
	if err != nil {
		die(w, r, err)
		return
	}

	imgURLs, err := c.FetchImageURLs(r.Context())
	if err != nil {
		die(w, r, err)
		return
	}

	// Pages are zero-indexed
	if page >= len(imgURLs) {
		err = fmt.Errorf("page %d is out of bounds, max %d", page, len(imgURLs)-1)
		die(w, r, err)
		return
	}

	imgURL := imgURLs[page]
	w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(imgURL.Path)))

	err = shared.QueryImage(r.Context(), imgURL, w)
	if err != nil {
		die(w, r, err)
	}
}

func epubHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	c, err := chapter.Fetch(r.Context(), id, r.URL.Query())
	if err != nil {
		die(w, r, err)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(".epub"))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.epub"`, c.FullTitle()))

	// TODO etags and caching?

	err = c.WriteEpub(r.Context(), w)
	if err != nil {
		die(w, r, err)
	}
}

func cbzHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
	}

	c, err := chapter.Fetch(r.Context(), id, r.URL.Query())
	if err != nil {
		die(w, r, err)
		return
	}

	// We're cheating a bit, the EPUBs that are created are also valid CBZs!
	// So as long as we tell the client that everything should work!
	w.Header().Set("Content-Type", mime.TypeByExtension(".cbz")) // TODO wrong mimetype?
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.cbz"`, c.FullTitle()))

	// TODO etags and caching?

	err = c.WriteEpub(r.Context(), w)
	if err != nil {
		die(w, r, err)
	}
}
