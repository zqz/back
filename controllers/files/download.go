package files

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zqzca/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/file"
)

// Download sends the entire file to the client.
func Download(c echo.Context) error {
	slug := c.Param("slug")
	f, err := file.FindBySlug(db.Connection, slug)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	// Build Etag
	etag := f.Hash
	res := c.Response()
	res.Header().Set("Content-Type", f.Type)
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

	// Build file and spit it out.
	fb := file.NewBuilder(db.Connection, f)

	if fb == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)

	fb.Copy(c.Response(), func() {
		// Once finished, flush output
		res.(http.Flusher).Flush()
	})

	return nil // Don't return anything, we sent data already.
}
