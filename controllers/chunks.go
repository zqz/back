package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/chunk"
	"github.com/zqzca/back/models/file"
)

func ChunkCreate(c *echo.Context) error {
	position, err := strconv.Atoi(c.Form("position"))
	if err != nil {
		fmt.Println("Failed to convert position")
		return c.NoContent(http.StatusBadRequest)
	}

	fileID := c.Form("file_id")

	tx := StartTransaction()

	// Make sure file exists.
	f, err := file.FindByID(tx, fileID)
	if err != nil {
		tx.Rollback()
		fmt.Println("File does not exist")
		return c.NoContent(http.StatusNotFound)
	}

	// Make sure chunk is not previously uploaded.
	if chunk.HaveChunkForFile(tx, fileID, position) {
		tx.Rollback()
		fmt.Println("Already received chunk at position", position)
		return c.NoContent(http.StatusNotAcceptable)
	}

	// Make sure we're only uploading chunks to incomplete files.
	if f.State != file.Incomplete {
		tx.Rollback()
		fmt.Println("File is not in incomplete state", f.State)
		return c.NoContent(http.StatusNotAcceptable)
	}

	req := c.Request()
	req.ParseMultipartForm(16 << 20)

	files := req.MultipartForm.File["data"]
	file := files[0]

	// Source file
	src, err := file.Open()
	if err != nil {
		fmt.Println("createChunk error: failed to open file", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer src.Close()

	var hash string
	if hash, err = lib.Hash(src); err != nil {
		fmt.Println("createChunk error: failed to hash file", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Destination file
	dstPath := fmt.Sprintf("files/chunks/%s", hash)

	var size int
	if size, err = storeChunk(src, dstPath); err != nil {
		fmt.Println("not sure:", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	chnk := &chunk.Chunk{
		FileID:   f.ID,
		Size:     size,
		Hash:     hash,
		Position: position,
	}

	chnk.Create(tx)
	chunks, _ := chunk.FindByFileID(tx, fileID)

	fmt.Println("i have :", len(*chunks), "chunks")
	fmt.Println("i need :", f.Chunks, "chunks")

	if len(*chunks) == f.Chunks {
		f.Process(tx)
	}

	fmt.Println("h: ", chnk.Hash)
	fmt.Println("f: ", fileID)
	fmt.Println("pos: ", position)
	tx.Commit()

	return c.NoContent(http.StatusOK)
}

func storeChunk(src io.Reader, path string) (int, error) {
	// Destination file
	dst, err := os.Create(path)

	if err != nil {
		return 0, err
	}

	defer dst.Close()

	var fileSize int64

	if fileSize, err = io.Copy(dst, src); err != nil {
		return int(fileSize), err
	}

	return int(fileSize), nil
}
