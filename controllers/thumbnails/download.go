package thumbnails

import (
	"github.com/zqzca/back/lib"
	"github.com/zqzca/echo"
)

// Download a file
func (t ThumbnailsController) Download(e echo.Context) error {
	fileID := e.Param("id")

	t, err := thumbnail.FindByID(t.DB, fileID)
	if err != nil {
		return err
	}

	path := lib.LocalPath(t.Hash)
	return e.File(path)
}
