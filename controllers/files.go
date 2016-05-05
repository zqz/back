package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

func FileIndex(c *echo.Context) error {
	var page uint = 0
	var per_page uint = 10

	files, err := models.FilePagination(page, per_page)

	if err != nil {
		fmt.Println("failed to fetch page:", err)
		return err
	}

	return c.JSON(http.StatusOK, files)
}

func FileDownload(c *echo.Context) error {
	file_id := c.Param("file_id")
	f, err := models.FileFindByID(file_id)

	// File must exist
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	// Build Etag
	etag := f.Hash
	res := c.Response()
	res.Header().Set(echo.ContentType, f.Type)
	res.Header().Set("Etag", etag)
	res.Header().Set("Cache-Control", "max-age=2592000") // 30 days

	// If set just return early.
	if match := c.Request().Header.Get("If-None-Match"); match != "" {
		fmt.Println("existing", match)
		fmt.Println("want", etag)
		if strings.Contains(match, etag) {
			return c.NoContent(http.StatusNotModified)
		}
	}

	// Build file and spit it out.
	fb := models.NewFileBuilder(f)

	if fb == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)

	fb.Copy(c.Response(), func() {
		res.Flush()
	})

	return nil
}

func FileCreate(c *echo.Context) error {
	f := &models.File{}

	if err := c.Bind(f); err != nil {
		return err
	}

	f.State = models.Incomplete
	f.Save()

	return c.JSON(http.StatusOK, f)
}

func FileStatus(c *echo.Context) error {
	hash := GetParam(c, "hash")
	f, err := models.FileFindByHash(hash)
	fs := models.FileStatusForFile(f)

	if err == nil {
		return c.JSON(http.StatusOK, fs)
	} else {
		fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}
}
