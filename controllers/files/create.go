package files

import (
	"net/http"

	"github.com/pressly/chi/render"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"

	"github.com/vattle/sqlboiler/boil"
)

func fileExistsWithHash(ex boil.Executor, hash string) (bool, error) {
	// Todo write an exists? for this
	count, err := models.Files(ex, qm.Where("hash=$1", hash)).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Create creates a file container in the database.
func (f Controller) Create(w http.ResponseWriter, r *http.Request) {
	file := &models.File{}

	if err := render.Bind(r.Body, file); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	if file.NumChunks < 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	exists, err := fileExistsWithHash(f.DB, file.Hash)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if exists {
		f.Debug("file exists with hash", "hash", file.Hash)
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, map[string]string{"error": "File Exists"})
		return
	}

	f.Debug("file doesnt exist with hash", "hash", file.Hash)
	file.State = lib.FileIncomplete

	if err := file.Insert(f.DB); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, file)
	return
}
