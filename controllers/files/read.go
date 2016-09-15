package files

import (
	"net/http"

	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	. "github.com/vattle/sqlboiler/queries/qm"
)

// Read returns a JSON payload to the client
func (f FileController) Read(e echo.Context) error {
	slug := e.Param("slug")
	file, err := models.Files(f.DB, Where("slug=$1", slug)).One()

	if err != nil {
		return e.NoContent(http.StatusNotFound)
	}

	return e.JSON(http.StatusOK, file)
}
