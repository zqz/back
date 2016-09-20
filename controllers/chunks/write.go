package chunks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/processors"
	"github.com/zqzca/echo"

	"github.com/vattle/sqlboiler/queries/qm"
)

const maxChunkSize = 5 * 1024 * 1024

func (c ChunkController) Write(e echo.Context) error {
	req := e.Request()
	contentLength := req.ContentLength()

	if contentLength == 0 {
		c.Debug("No content length")
		return e.NoContent(http.StatusLengthRequired)
	}

	if contentLength > maxChunkSize {
		c.Debug("Content length > maxChunkSize", "length", contentLength, "max", maxChunkSize)
		return e.NoContent(http.StatusRequestEntityTooLarge)
	}

	// Make sure the file exists.
	fid := e.Param("file_id")
	f, err := models.FindFile(c.DB, fid)
	if err != nil {
		c.Debug("File not found", "hash", fid)
		return e.NoContent(http.StatusNotFound)
	}

	chunkID, err := strconv.Atoi(e.Param("chunk_id"))
	if err != nil {
		c.Error("Failed to convert chunk id into integer")
		return err
	}

	clientHash := e.Param("hash")
	if c.chunkExists(f.ID, clientHash) {
		c.Warn(
			"Chunk Already exists",
			"ID", clientHash,
		)

		return e.NoContent(http.StatusConflict)
	}

	// Actually read file...
	buf, err := ioutil.ReadAll(req.Body())
	if err != nil {
		c.Error("Failed to read Request Body", "error", err)

		return e.NoContent(http.StatusBadRequest)
	}
	time.Sleep(500 * time.Millisecond)

	b := bytes.NewReader(buf)
	hash, _ := lib.Hash(b)

	// Does this need to be handled? Can a bytesbuffer error on seek?
	b.Seek(0, os.SEEK_SET)

	c.Debug(
		"Chunk Received",
		"Request Size", contentLength,
		"Size", b.Size(),
		"Hash", hash,
	)

	if hash != clientHash {
		c.Warn(
			"Hash does not match what client specified",
			"Client", clientHash,
			"Server", hash,
		)
		return e.NoContent(422) // Unprocessable Entity
	}

	// Destination file
	dstPath := filepath.Join("files", "chunks", hash)

	var size int
	if size, err = c.storeChunk(b, dstPath); err != nil {
		c.Error("Failed to store chunk", "Error", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	chunk := &models.Chunk{
		FileID:   f.ID,
		Position: int(chunkID),
		Size:     size,
		Hash:     hash,
	}

	if err = chunk.Insert(c.DB); err != nil {
		c.Error("Failed to insert chunk in DB", "Error", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	c.checkFinished(f)

	return e.NoContent(http.StatusCreated)
}

func (c ChunkController) chunkExists(fid string, hash string) bool {
	chunkCount, err := models.Chunks(c.DB, qm.Where("file_id=$1 and hash=$2", fid, hash)).Count()

	if err != nil {
		c.Error("Failed to look up chunk count", err)
		return false
	}

	return chunkCount > 0
}

// storeChunk writes the chunk data from src to a new file at path.
func (c ChunkController) storeChunk(src io.Reader, path string) (int, error) {
	dst, err := os.Create(path)

	if err != nil {
		return 0, err
	}

	defer dst.Close()

	var fileSize int64

	if fileSize, err = io.Copy(dst, src); err != nil {
		c.Error(
			"Failed to copy chunk data to destination",
			"Destination", path,
			"Error", err,
		)

		return int(fileSize), err
	}

	return 0, nil
}

func (c ChunkController) checkFinished(f *models.File) {
	chunks, err := models.Chunks(c.DB, qm.Where("file_id=$1", f.ID)).All()

	if err != nil {
		c.Error("Failed to lookup chunks", "Error", err)
		return
	}

	completedChunks := len(chunks)
	requiredChunks := f.NumChunks

	fmt.Println("Completed Chunks:", completedChunks)
	fmt.Println("Required:", requiredChunks)

	if completedChunks != int(requiredChunks) {
		c.Info(
			"File not finished",
			"Received", completedChunks,
			"Total", requiredChunks,
		)

		return
	}

	go func() {
		err = processors.CompleteFile(c.Dependencies, f)

		if err != nil {
			c.Error("Failed to finish file", "error", err, "name", f.Name, "id", f.ID)
			return
		}
		c.Info("Finished File", "name", f.Name, "id", f.ID)
	}()
}
