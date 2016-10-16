package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Download sends the entire file to the client.
func (f Controller) Download(e echo.Context) error {
	slug := e.Param("slug")
	file, err := models.Files(f.DB, qm.Where("slug=$1", slug)).One()
	if err != nil {
		return e.NoContent(http.StatusNotFound)
	}

	// Build Etag
	etag := file.Hash
	res := e.Response()
	res.Header().Set("Content-Type", file.Type)
	res.Header().Set("Etag", etag)
	res.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	disposition := fmt.Sprintf("inline; filename=%s", file.Name)
	res.Header().Set("Content-Disposition", disposition)

	// If set just return early.
	if match := e.Request().Header().Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			go lib.TrackDownload(f.DB, file.ID, e, true)
			return e.NoContent(http.StatusNotModified)
		}
	}

	data, err := os.Open(lib.LocalPath(file.Hash))
	if err != nil {
		return e.NoContent(http.StatusNotFound)
	}

	go lib.TrackDownload(f.DB, file.ID, e, false)

	defer data.Close()
	_, err = io.Copy(res, data)
	return err
}
