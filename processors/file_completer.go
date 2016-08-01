package processors

import (
	"github.com/pkg/errors"
	"github.com/zqzca/back/models"
)

// CompleteFile builds the file from chunks and then generates thumbnails
func CompleteFile(f *models.File) error {
	reader, err := BuildFile(f)
	if err != nil {
		return errors.Wrap(err, "Failed to complete building file")
	}

	thumbHash, thumbSize, err := CreateThumbnail(reader)
	if err != nil {
		return errors.Wrap(err, "Failed to create thumbnail")
	}

	t := models.Thumbnail{
		Hash:   thumbHash,
		Size:   thumbSize,
		FileID: f.ID,
	}

	t.Insert()

	return nil
}
