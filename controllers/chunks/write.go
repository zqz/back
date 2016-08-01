package chunks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	. "github.com/nullbio/sqlboiler/boil/qm"
)

func Write(e echo.Context) error {
	req := e.Request()
	length := req.ContentLength()

	if length == 0 {
		return e.NoContent(http.StatusLengthRequired)
	}

	if length > 5*1024*1024 {
		return e.NoContent(http.StatusRequestEntityTooLarge)
	}

	fileID := e.Param("file_id")
	if !fileExists(fileID) {
		return e.NoContent(http.StatusNotFound)
	}

	clientHash := e.Param("hash")
	if chunkExists(fileID, clientHash) {
		return e.NoContent(http.StatusConflict)
	}

	// Actually read file.
	buf, err := ioutil.ReadAll(req.Body())

	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}

	b := bytes.NewReader(buf)
	hash, _ := lib.Hash(b)

	_, err := b.Seek(0, io.SET_SEEK)

	fmt.Println("Length: ", length)
	fmt.Println("Size:", b.Size())
	fmt.Println("Hash:", hash)

	if hash != clientHash {
		return e.NoContent(422) // Unprocessable Entity
	}

	// Destination file
	dstPath := filepath.Join("files", "chunks", hash)

	var size int
	if size, err = storeChunk(b, dstPath); err != nil {
		fmt.Println("failed to store chunk:", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	c := &models.Chunk{
		FileID:   fileID,
		Position: int(chunkID),
		Size:     size,
		Hash:     hash,
	}

	err = c.Insert()
	if err != nil {
		return e.NoContent(http.StatusInternalServerError)
	}

	go checkFinished(fileID)

	return e.NoContent(http.StatusCreated)
}

func fileExists(fid string) bool {
	return models.Files(Where("file_id=$1", fid)).Count() > 0
}

func chunkExists(fid string, hash string) bool {
	return models.Chunks(Where("file_id=$1 and hash=$2", fid, hash)).Count() > 0
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
	chunks, err := models.Chunks(Where("file_id=$1", fid)).All()

	if err != nil {
		fmt.Println("Failed to query for all chunks with file id:", fid)
		return
	}

	f, err := models.FileFind(fid)

	if err != nil {
		fmt.Println("Failed to query for file", err)
		return
	}

	completed_chunks := len(chunks)
	required_chunks := f.Chunks

	if completed_chunks == required_chunks {
		err = FileCompleter(f)

		if err != nil {
			return errors.Wrap(err, "Failed to complete file")
		}
		// fmt.Println("processing")
		// tx := Begin()
		// // f.Process(tx)
		// tx.Commit()
	}
}
