package files

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"
)

//Index returns a list of files
func (f Controller) Index(e echo.Context) error {
	pageStr := e.QueryParam("page")
	perPageStr := e.QueryParam("per_page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		f.Warn("Failed to decode page param", "e", err)
		return e.NoContent(http.StatusBadRequest)
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		f.Warn("Failed to decode per_page param")
		return e.NoContent(http.StatusBadRequest)
	}

	files, err := models.Files(f.DB, qm.Limit(perPage), qm.Offset(page*perPage)).All()
	if err != nil {
		fmt.Println("failed to fetch page:", err)
		return err
	}

	return e.JSON(http.StatusOK, files)
}
