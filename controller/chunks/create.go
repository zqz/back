package chunks

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zqzca/back/models"
	"github.com/zqzca/back/processors"
	"github.com/zqzca/back/serializer"

	"github.com/vattle/sqlboiler/queries/qm"
)

// Create receives chunk data via a POST.
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
		Position: u.chunkID,
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

	if completedChunks != requiredChunks {
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

		err = processors.CompleteFile(c.Dependencies, *f)

		if err != nil {
			c.Error("Failed to finish file", "error", err, "name", f.Name, "id", f.ID)
			return
		}

		if len(wsID) > 0 {
			c.Info("Sending WS msg", "ws", wsID)
			c.WS.WriteClient(wsID, "file:completed", f)
			de := serializer.NewDashboardItemFromFile(c.DB, f)
			c.WS.Broadcast("file:added", de)
		} else {
			c.Info("No WS ID")
		}

		c.Info("Finished File", "name", f.Name, "id", f.ID)
	}()
}

func parseRequest(r *http.Request) *upload {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("err", err.Error())
		return nil
	}

	positionStr := m["position"][0]
	position, err := strconv.Atoi(positionStr)
	if err != nil {
		fmt.Println("err", err.Error())
		return nil
	}

	return &upload{
		request:    r,
		size:       int(r.ContentLength),
		chunkID:    position,
		fileID:     m["file_id"][0],
		remoteHash: m["hash"][0],
		wsID:       m["ws_id"][0],
	}
}

func (c Controller) storeWebsocket(fID string, ws string) {
	c.wsFileIDsLock.Lock()
	c.Info("Storing WS for File", "ws", ws, "file", fID)
	c.wsFileIDs[fID] = ws
	c.wsFileIDsLock.Unlock()
}
