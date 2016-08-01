package chunks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/chunk"
	"github.com/zqzca/back/models/file"
	"github.com/zqzca/echo"
)

func Write(c echo.Context) error {
	req := c.Request()
	length := req.ContentLength()

	if length == 0 {
		return c.NoContent(http.StatusLengthRequired)
	}

	if length > 5*1024*1024 {
		return c.NoContent(http.StatusRequestEntityTooLarge)
	}

	fileID := c.Param("file_id")
	if !fileExists(tx, fileID) {
		return c.NoContent(http.StatusNotFound)
	}

	clientHash := c.Param("hash")
	if chunkExists(tx, fileID, clientHash) {
		return c.NoContent(http.StatusConflict)
	}

	// Actually read file.
	buf, err := ioutil.ReadAll(req.Body())

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	b := bytes.NewReader(buf)
	hash, _ := lib.Hash(b)

	_, err := b.Seek(0, io.SET_SEEK)

	fmt.Println("Length: ", length)
	fmt.Println("Size:", b.Size())
	fmt.Println("Hash:", hash)

	if hash != clientHash {
		return c.NoContent(422) // Unprocessable Entity
	}

	tx := db.StartTransaction()
	// Destination file
	dstPath := filepath.Join("files", "chunks", hash)

	var size int
	if size, err = storeChunk(b, dstPath); err != nil {
		fmt.Println("failed to store chunk:", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	chnk := &chunk.Chunk{
		FileID:   fileID,
		Position: int(chunkID),
		Size:     size,
		Hash:     hash,
	}

	err = chnk.Create(tx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	tx.Commit()

	go checkFinished(fileID)

	return c.NoContent(http.StatusCreated)
}

func fileExists(ex db.Executor, fid string) bool {
	_, err := file.FindByID(ex, fid)

	return err == nil
}

func chunkExists(ex db.Executor, fid string, hash string) bool {
	return chunk.HaveChunkForFileWithHash(ex, fid, hash)
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

func checkFinished(fid string) {
	chunks, err := chunk.FindByFileID(db.Connection, fid)

	if err != nil {
		fmt.Println("Failed to query for all chunks with file id:", fid)
		return
	}

	f, err := file.FindByID(db.Connection, fid)

	if err != nil {
		fmt.Println("Failed to query for file", err)
		return
	}

	completed_chunks := len(*chunks)
	required_chunks := f.Chunks

	if completed_chunks == required_chunks {
		fmt.Println("processing")
		tx := db.StartTransaction()
		f.Process(tx)
		tx.Commit()
	}
}
