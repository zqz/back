package thumbnails

import (
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"
)

// Download a file
func (t ThumbnailsController) Download(e echo.Context) error {
	id := e.Param("id")

	thumb, err := models.ThumbnailFind(t.DB, id)
	if err != nil {
		return err
	}

	return e.File(lib.LocalPath(thumb.Hash))
}
