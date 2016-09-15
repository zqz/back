package thumbnails

import (
	"net/http"

	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"
)

// Download a file
func (t ThumbnailsController) Download(e echo.Context) error {
	id := e.Param("id")

	if len(id) != 36 {
		return e.NoContent(http.StatusNotFound)
	}

	thumb, err := models.FindThumbnail(t.DB, id)
	if err != nil {
		return err
	}

	return e.File(lib.LocalPath(thumb.Hash))
}
