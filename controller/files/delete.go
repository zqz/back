package files

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/zqzca/back/models"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Delete removes a file, it's data, chunks and any thumbnails.
func (f Controller) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if len(slug) == 0 {
		http.Error(w, "No slug specified", http.StatusBadRequest)
		return
	}

	tx, err := f.DB.Begin()
	if err != nil {
		f.Error("Failed to create transaction")
		http.Error(w, "Failed to DB create transaction", 500)
		return
	}

	file, err := models.Files(tx, qm.Where("slug=?", slug)).One()
	if err != nil {
		f.Error("Failed to fetch file")
		http.Error(w, "Failed to find File in DB", 404)
		return
	}

	err = file.Chunks(tx).DeleteAll()
	if err != nil {
		f.Error("Failed to delete chunks", "err", err.Error())
		_ = tx.Rollback()
		http.Error(w, "Failed to delete chunks", 500)
		return
	}

	err = file.Thumbnails(tx).DeleteAll()
	if err != nil {
		f.Error("Failed to delete thumbnails")
		_ = tx.Rollback()
		http.Error(w, "Failed to delete thumbnails", 500)
		return
	}

	if err = file.Delete(tx); err != nil {
		f.Error("Failed to delete file record")
		http.Error(w, "Failed to delete file from DB", 500)
		return
	}

	if err = tx.Commit(); err != nil {
		f.Error("Failed to commit transaction")
		http.Error(w, "Failed to commit TX", 500)
		return
	}

	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	return
}
