package files

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	. "github.com/vattle/sqlboiler/boil/qm"
)

// Download sends the entire file to the client.
func (f FileController) Download(e echo.Context) error {
	slug := e.Param("slug")
	file, err := models.Files(f.DB, Where("slug=$1", slug)).One()
	if err != nil {
		return e.NoContent(http.StatusNotFound)
	}

	// Build Etag
	etag := file.Hash
	res := e.Response()
	res.Header().Set("Content-Type", file.Type)
	res.Header().Set("Etag", etag)
	res.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	res.Header().Set("Content-Disposition", "inline")    // 30 days

	// If set just return early.
	if match := e.Request().Header().Get("If-None-Match"); match != "" {
		f.Debug("existing", "match", match)
		f.Debug("want", "etag", etag)
		if strings.Contains(match, etag) {
			return e.NoContent(http.StatusNotModified)
		}
	}

	data, err := os.Open(lib.LocalPath(file.Hash))
	if err != nil {
		return e.NoContent(http.StatusNotFound)
	}
	defer data.Close()
	_, err = io.Copy(res, data)
	return err
}
