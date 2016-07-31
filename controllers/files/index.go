package files

import (
	"fmt"
	"net/http"

	"github.com/zqzca/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/file"
)

//Index returns a list of files
func Index(c echo.Context) error {
	page := 0
	perPage := 10

	files, err := file.Pagination(db.Connection, page, perPage)

	if err != nil {
		fmt.Println("failed to fetch page:", err)
		return err
	}

	return c.JSON(http.StatusOK, files)
}
