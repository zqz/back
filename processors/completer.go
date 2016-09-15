package processors

import (
	"fmt"

	"github.com/pkg/errors"
	. "github.com/vattle/sqlboiler/queries/qm"
	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
)

// CompleteFile builds the file from chunks and then generates thumbnails
func CompleteFile(deps controllers.Dependencies, f *models.File) error {
	fmt.Println("Completing File", f.Name, f.ID)
	tx, err := deps.DB.Begin()

	if f == nil {
		fmt.Println("Fucking goof, it dont work if its nil")
	}

	f.Reload(tx)

	if f.State == lib.FileProcessing {
		fmt.Println("Something else already processing the file")
		return nil
	}

	models.Thumbnails(tx, Where("file_id=$1", f.ID)).DeleteAll()

	f.State = lib.FileProcessing
	f.Update(tx, "state")
	deps.Debug("Locking", "ID", f.ID)

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

	if len(thumbHash) > 0 {
		t := models.Thumbnail{
			Hash:   thumbHash,
			Size:   thumbSize,
			FileID: f.ID,
		}

		t.Insert(deps.DB)
	}

	f.State = lib.FileFinished
	f.Update(tx, "state")
	err = tx.Commit()

	if err != nil {
		fmt.Println("failed to commit tx", f.ID)
	}

	err = Cleanup(deps, f)

	if err != nil {
		fmt.Println("Failed to cleanup file", f.ID)
	}

	return nil
}
