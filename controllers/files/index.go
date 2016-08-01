package files

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	. "github.com/nullbio/sqlboiler/boil/qm"
)

//Index returns a list of files
func (f FileController) Index(e echo.Context) error {
	pageStr := e.Param("page")
	perPageStr := e.Param("per_page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}

	files, err := models.Files(f.DB, Limit(perPage), Offset(page*perPage)).All()
	if err != nil {
		fmt.Println("failed to fetch page:", err)
		return err
	}

	return e.JSON(http.StatusOK, files)
}
