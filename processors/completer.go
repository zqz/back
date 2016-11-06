package processors

import (
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/zqzca/back/dependencies"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
)

// CompleteFile builds the file from chunks and then generates thumbnails
func CompleteFile(deps dependencies.Dependencies, f *models.File) error {
	deps.Info("Processing File", "name", f.Name, "id", f.ID)

	tx, err := deps.DB.Begin()
	if err != nil {
		deps.Error("Failed to create transaction")
		return err
	}

	if err = f.Reload(tx); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Failed to reload the file")
	}

	if f.State == lib.FileProcessing {
		tx.Rollback()
		return errors.Wrap(err, "This file is already being processed")
	}

	// Delete all thumbnails
	err = models.Thumbnails(tx, qm.Where("file_id=?", f.ID)).DeleteAll()
	if err != nil {
		deps.Info("No previous thumnails")
	}

	f.State = lib.FileProcessing
	if err = f.Update(tx, "state"); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Failed to update state")
	}

	reader, err := BuildFile(deps, f)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Failed to complete building file")
	}

	thumbHash, thumbSize, err := CreateThumbnail(deps, reader)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Failed to create thumbnail")
	}

	if len(thumbHash) == 0 {
		tx.Rollback()
		return errors.Wrap(err, "No thumbnail created")
	}

	t := models.Thumbnail{
		Hash:   thumbHash,
		Size:   thumbSize,
		FileID: f.ID,
	}

	if err = t.Insert(tx); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Failed to insert Thumbnail")
	}

	f.State = lib.FileFinished
	if err = f.Update(tx, "state"); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Failed to set state")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "Failed to commit transaction")
	}

	if err = Cleanup(deps, f); err != nil {
		return errors.Wrap(err, "Failed to cleanup file")
	}

	deps.Info("Processed File", "name", f.Name, "id", f.ID)
	return nil
}
