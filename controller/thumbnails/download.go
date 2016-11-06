package thumbnails

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
)

// Download a file
func (t Controller) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if len(id) != 36 {
		http.Error(w, "Thumbnail not found", 404)
		return
	}

	thumb, err := models.FindThumbnail(t.DB, id)
	if err != nil {
		http.Error(w, "Thumbnail not found", 404)
		return
	}

	http.ServeFile(w, r, lib.LocalPath(thumb.Hash))
}
