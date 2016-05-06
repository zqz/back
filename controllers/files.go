package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models/file"
)

func FileIndex(c *echo.Context) error {
	tx := StartTransaction()
	defer tx.Rollback()

	page := 0
	perPage := 10

	files, err := file.Pagination(tx, page, perPage)

	if err != nil {
		fmt.Println("failed to fetch page:", err)
		return err
	}

	return c.JSON(http.StatusOK, files)
}

func FileDownload(c *echo.Context) error {
	tx := StartTransaction()
	defer tx.Rollback()

	fileID := c.Param("file_id")
	f, err := file.FindByID(tx, fileID)
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
	fb := file.NewBuilder(tx, f)

	if fb == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)

	fb.Copy(c.Response(), func() {
		res.Flush()
	})

	return nil // Don't return anything, we sent data already.
}

func FileCreate(c *echo.Context) error {
	f := &file.File{}

	if err := c.Bind(f); err != nil {
		return err
	}

	tx := StartTransaction()
	f.State = file.Incomplete
	f.Create(tx)
	tx.Commit()

	return c.JSON(http.StatusOK, f)
}

func FileStatus(c *echo.Context) error {
	tx := StartTransaction()
	defer tx.Rollback()
	hash := GetParam(c, "hash")
	f, err := file.FindByHash(tx, hash)
	fs := file.StatusForFile(tx, f)

	if err == nil {
		return c.JSON(http.StatusOK, fs)
	} else {
		fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}
}
