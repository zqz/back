package files

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/zqzca/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/file"
)

// Read returns a JSON payload to the client
func Read(c echo.Context) error {
	slug := c.Param("slug")
	f, err := file.FindBySlug(db.Connection, slug)

	spew.Dump(f)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, f)
}
