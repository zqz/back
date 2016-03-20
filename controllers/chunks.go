package controllers

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

func ChunkCreate(c *echo.Context) error {
	position, err := strconv.Atoi(c.Form("position"))
	if err != nil {
		fmt.Println("Failed to convert position")
		return c.NoContent(http.StatusBadRequest)
	}

	fileID := c.Form("file_id")

	// Make sure file exists.
	f, err := models.FileFindByID(fileID)
	if err != nil {
		fmt.Println("File does not exist")
		return c.NoContent(http.StatusNotFound)
	}

	// Make sure chunk is not previously uploaded.
	if _, err = models.ChunkFindByFileIDAndPosition(fileID, position); err == nil {
		fmt.Println("Already received chunk at position", position)
		return c.NoContent(http.StatusNotAcceptable)
	}

	// Make sure we're only uploading chunks to incomplete files.
	if f.State != models.Incomplete {
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
		return err
	}
	defer src.Close()

	var hash string
	if hash, err = hashFile(src); err != nil {
		return err
	}

	// Make sure chunk is not previously uploaded.
	if _, err = models.ChunkFindByHash(hash); err == nil {
		fmt.Println("Already received hash: ", hash)
		return c.NoContent(http.StatusNotAcceptable)
	}

	// Destination file
	dstPath := fmt.Sprintf("files/chunks/%s", hash)

	var size int
	if size, err = storeChunk(src, dstPath); err != nil {
		return err
	}

	chunk := &models.Chunk{
		FileID:   f.ID,
		Size:     size,
		Hash:     hash,
		Position: position,
	}

	chunk.Save()
	chunks, _ := models.ChunksByFileID(fileID)

	fmt.Println("i have :", chunks, "chunks")
	fmt.Println("i need :", f.Chunks, "chunks")
	if len(chunks) == f.Chunks {
		f.Process()
	}

	fmt.Println(hash)
	fmt.Println("h: ", hash)
	fmt.Println("p: ", dstPath)
	fmt.Println("f: ", fileID)
	fmt.Println("pos: ", position)

	return c.NoContent(http.StatusOK)
}

func hashFile(src io.ReadSeeker) (string, error) {
	h := sha1.New()

	if _, err := io.Copy(h, src); err != nil {
		fmt.Println(err)
		return "", err
	}

	if _, err := src.Seek(0, os.SEEK_SET); err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
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
