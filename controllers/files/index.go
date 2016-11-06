package files

import (
	"net/http"

	"github.com/pressly/chi/render"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/zqzca/back/models"
)

//Index returns a list of files
func (c Controller) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	perPage := ctx.Value(1001).(int)
	page := ctx.Value(1001).(int)

	files, err := models.Files(
		c.DB,
		qm.OrderBy("created_at desc"),
		qm.Limit(perPage),
		qm.Offset(page*perPage),
	).All()

	if err != nil {
		http.Error(w, "Failed to fetch page", 500)
		return
	}

	render.JSON(w, r, files)
}
