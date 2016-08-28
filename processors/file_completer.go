package processors

import (
	"github.com/pkg/errors"
	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/models"
)

// CompleteFile builds the file from chunks and then generates thumbnails
func CompleteFile(deps controllers.Dependencies, f *models.File) error {
	reader, err := BuildFile(deps, f)
	if err != nil {
		return errors.Wrap(err, "Failed to complete building file")
	}

	thumbHash, thumbSize, err := CreateThumbnail(deps, reader)
	if err != nil {
		return errors.Wrap(err, "Failed to create thumbnail")
	}

	if thumbSize == 0 {
		return nil
	}

	t := models.Thumbnail{
		Hash:   thumbHash,
		Size:   thumbSize,
		FileID: f.ID,
	}

	t.Insert(deps.DB)

	return nil
}
