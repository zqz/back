package files

import (
	"net/http"

	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/processors"
	"github.com/labstack/echo"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Process builds thumbnails
func (f Controller) Process(e echo.Context) error {
	f.Debug("processing")

	slug := e.Param("slug")
	file, err := models.Files(f.DB, qm.Where("slug=$1", slug)).One()

	if err != nil {
		return err
	}

	if file.State == lib.FileProcessing {
		f.Debug("file already being processed")
		return e.NoContent(http.StatusConflict)
	}

	err = processors.CompleteFile(f.Dependencies, file)

	if err != nil {
		f.Debug("failed to complete file", "err", err)
		return err
	}

	f.Info("Finished File", "name", file.Name, "id", file.ID)

	return e.NoContent(http.StatusOK)
}
