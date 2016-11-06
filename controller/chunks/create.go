package chunks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/processors"

	"github.com/vattle/sqlboiler/queries/qm"
)

const maxChunkSize = 5 * 1024 * 1024

type chunkUpload struct {
	chunkID    int
	fileID     string
	size       int
	wsID       string
	remoteHash string
	localHash  string
	data       []byte

	request *http.Request
}

func parseRequest(r *http.Request) *chunkUpload {
	c := chunkUpload{}

	u := r.URL
	fmt.Println(u.RawQuery)

	m, err := url.ParseQuery(u.RawQuery)

	if err != nil {
		fmt.Println("err", err.Error())
	}

	c.request = r
	chunkIDStr := m["position"][0]
	c.chunkID = -1
	c.chunkID, _ = strconv.Atoi(chunkIDStr)
	c.size = int(r.ContentLength)
	c.fileID = m["file_id"][0]
	c.remoteHash = m["hash"][0]
	c.wsID = m["ws_id"][0]

	return &c
}

func (u *chunkUpload) validData() (bool, error) {
	if u.remoteHash != u.localHash {
		return false, errors.New("Hash does not match")
	}

	if u.size != len(u.data) {
		return false, errors.New("Incorrect size given")
	}

	return true, nil
}

func (u *chunkUpload) validRequest() (bool, error) {
	if u.size == 0 {
		return false, errors.New("Chunk has no size")
	}

	if u.size > maxChunkSize {
		return false, errors.New("Chunk is too big")
	}

	if len(u.fileID) == 0 {
		return false, errors.New("No FileID specified")
	}

	if u.chunkID < 0 {
		return false, errors.New("No ChunkID specified")
	}

	if len(u.remoteHash) == 0 {
		return false, errors.New("No Hash specified")
	}

	return true, nil
}

func (u *chunkUpload) loadData() error {
	var err error
	u.data, err = ioutil.ReadAll(u.request.Body)
	return err
}

func (u *chunkUpload) hashData() error {
	br := bytes.NewReader(u.data)
	var err error
	u.localHash, err = lib.Hash(br)
	br.Seek(0, os.SEEK_SET)
	return err
}

func (u *chunkUpload) storeData() error {
	dst, err := os.Create(u.localPath())
	if err != nil {
		return err
	}
	defer dst.Close()

	br := bytes.NewReader(u.data)
	if _, err := io.Copy(dst, br); err != nil {
		return errors.Wrap(err, "Failed to copy chunk data to dst")
	}

	return nil
}

func (u *chunkUpload) localPath() string {
	return filepath.Join("files", "chunks", u.localHash)
}

func (c Controller) storeWebsocket(fID string, ws string) {
	c.wsFileIDsLock.Lock()
	c.Info("Storing WS for File", "ws", ws, "file", fID)
	c.wsFileIDs[fID] = ws
	c.wsFileIDsLock.Unlock()
}

func (c Controller) Create(w http.ResponseWriter, r *http.Request) {
	u := parseRequest(r)

	if ok, err := u.validRequest(); !ok {
		c.Debug("Invalid Request")
		http.Error(w, err.Error(), 400)
		return
	}

	f, err := models.FindFile(c.DB, u.fileID)
	if err != nil {
		c.Debug("File not found", "file_id", u.fileID)
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}

	if err := u.loadData(); err != nil {
		c.Debug("Failed to read chunk data")
		http.Error(w, err.Error(), 500)
		return
	}

	if err := u.hashData(); err != nil {
		c.Error("Failed to hash data")
		http.Error(w, "Failed to hash data", http.StatusInternalServerError)
		return
	}

	if ok, err := u.validData(); !ok {
		c.Error("Data inconsistency")
		http.Error(w, err.Error(), 422)
		return
	}

	if c.chunkExists(f.ID, u.localHash) {
		c.Warn("Chunk Already exists", "file_id", u.fileID, "chunk_id", u.chunkID)
		// TODO: check if file is finished
		http.Error(w, "Chunk Already exists", http.StatusConflict)
		return
	}

	// Debug remote this.
	time.Sleep(500 * time.Millisecond)

	c.Debug(
		"Chunk Received",
		"Request Size", u.size,
		"Size", len(u.data),
		"Hash", u.localHash,
	)

	if err := u.storeData(); err != nil {
		c.Error("Failed to store chunk", "Error", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	chunk := &models.Chunk{
		FileID:   u.fileID,
		Position: int(u.chunkID),
		Size:     u.size,
		Hash:     u.localHash,
	}

	if len(u.wsID) == 36 {
		c.storeWebsocket(u.fileID, u.wsID)
	}

	if err := chunk.Insert(c.DB); err != nil {
		c.Error("Failed to insert chunk in DB", "Error", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	c.checkFinished(f)
	return
}

func (c Controller) chunkExists(fid string, hash string) bool {
	chunkCount, err := models.Chunks(c.DB, qm.Where("file_id=$1 and hash=$2", fid, hash)).Count()

	if err != nil {
		c.Error("Failed to look up chunk count", err)
		return false
	}

	return chunkCount > 0
}

// storeChunk writes the chunk data from src to a new file at path.
func (c Controller) storeChunk(src io.Reader, path string) (int, error) {
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

func (c Controller) checkFinished(f *models.File) {
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
		c.wsFileIDsLock.RLock()
		wsID := c.wsFileIDs[f.ID]
		c.wsFileIDsLock.RUnlock()

		err = processors.CompleteFile(c.Dependencies, f)

		if err != nil {
			c.Error("Failed to finish file", "error", err, "name", f.Name, "id", f.ID)
			return
		}

		if len(wsID) > 0 {
			c.Info("Sending WS msg", "ws", wsID)
			c.WS.WriteClient(wsID, "file:completed", f)
		} else {
			c.Info("No WS ID")
		}

		c.Info("Finished File", "name", f.Name, "id", f.ID)
	}()
}
