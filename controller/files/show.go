package files

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/serializer"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Show returns a JSON payload to the client
func (c Controller) Show(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	f, err := models.Files(c.DB, qm.Where("slug=$1", slug)).One()
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}

	s := serializer.ForFile(c.DB, f)

	render.JSON(w, r, s)
}
