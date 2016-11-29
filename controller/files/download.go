package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Download sends the entire file to the client.
func (f Controller) Download(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if len(slug) == 0 {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "")
		return
	}

	file, err := models.Files(f.DB, qm.Where("slug=$1", slug)).One()
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
		return
	}

	// Build Etag
	etag := file.Hash
	w.Header().Set("Content-Type", file.Type)
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	disposition := fmt.Sprintf("inline; filename=%s", file.Name)
	w.Header().Set("Content-Disposition", disposition)

	// If set just return early.
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			go lib.TrackDownload(f.DB, file.ID, r, true)
			render.Status(r, http.StatusNotModified)
			render.PlainText(w, r, "")
			return
		}
	}

	data, err := os.Open(lib.LocalPath(file.Hash))
	if err != nil {
		render.Status(r, http.StatusNotModified)
		render.PlainText(w, r, "")
		return
	}
	defer data.Close()

	go lib.TrackDownload(f.DB, file.ID, r, false)

	if _, err := io.Copy(w, data); err != nil {
		http.Error(w, "Failed to write response", 500)
	}
}
