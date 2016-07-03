package chunks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/chunk"
	"github.com/zqzca/back/models/file"
)

func Write(c echo.Context) error {
	req := c.Request()
	length := req.ContentLength()

	if length > 5*1024*1024 {
		return c.NoContent(http.StatusRequestEntityTooLarge)
	}

	tx := db.StartTransaction()

	fileID := c.Param("file_id")
	chunkID, err := strconv.ParseInt(c.Param("chunk_id"), 10, 16)

	if err != nil {
		fmt.Println("Can not parse chunk id", c.Param("chunk_id"))
		return c.NoContent(http.StatusInternalServerError)
	}

	if !fileExists(tx, fileID) {
		return c.NoContent(http.StatusNotFound)
	}

	if chunkExists(tx, fileID, int(chunkID)) {
		return c.NoContent(http.StatusConflict)
	}

	// Actually read file.
	buf, _ := ioutil.ReadAll(req.Body())
	b := bytes.NewReader(buf)
	hash, _ := lib.Hash(b)

	fmt.Println("Length: ", length)
	fmt.Println("Size:", b.Size())
	fmt.Println("Hash:", hash)

	// Destination file
	dstPath := fmt.Sprintf("files/chunks/%s", hash)

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

	chnk.Create(tx)

	tx.Commit()

	go checkFinished(fileID)

	return c.NoContent(http.StatusCreated)
}

func fileExists(ex db.Executor, fid string) bool {
	_, err := file.FindByID(ex, fid)

	return err == nil
}

func chunkExists(ex db.Executor, fid string, cid int) bool {
	return chunk.HaveChunkForFile(ex, fid, cid)
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
