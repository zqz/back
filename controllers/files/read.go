package files

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/zqzca/back/models"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Read returns a JSON payload to the client
func (f Controller) Read(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	file, err := models.Files(f.DB, qm.Where("slug=$1", slug)).One()

	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}

	render.JSON(w, r, file)
}
