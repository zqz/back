package thumbnails

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/zqzca/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/thumbnail"
)

func Download(c echo.Context) error {
	fileID := c.Param("id")
	t, err := thumbnail.FindByID(db.Connection, fileID)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	// Build Etag
	etag := t.Hash
	res := c.Response()
	res.Header().Set("Content-Type", "image/jpeg")
	res.Header().Set("Etag", etag)
	res.Header().Set("Cache-Control", "max-age=2592000") // 30 days

	// If set just return early.
	if match := c.Request().Header().Get("If-None-Match"); match != "" {
		fmt.Println("existing", match)
		fmt.Println("want", etag)
		if strings.Contains(match, etag) {
			return c.NoContent(http.StatusNotModified)
		}
	}

	path := lib.LocalPath(t.Hash)
	f, err := os.Open(path)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)
	io.Copy(c.Response(), f)
	res.(http.Flusher).Flush()

	return nil // Don't return anything, we sent data already.
}
