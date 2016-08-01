package files

import (
	"fmt"
	"net/http"

	"github.com/nullbio/sqlboiler/boil"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	. "github.com/nullbio/sqlboiler/boil/qm"
)

type fileStatus struct {
	ID             string   `json:"id"`
	State          string   `json:"state"`
	ChunksReceived []string `json:"chunks_received,omitempty"`
	ChunksNeeded   int      `json:"chunks_needed,omitempty"`
}

func statusForFile(ex boil.Executor, f *models.File) *fileStatus {
	var state string

	switch f.State {
	case lib.FileIncomplete:
		state = "incomplete"
	case lib.FileProcessing:
		state = "processing"
	case lib.FileFinished:
		state = "finished"
	}

	chunksNeeded := 0
	chunks, err := models.Chunks(ex, Where("file_id=$1", f.ID)).All()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	numChunks := len(chunks)
	chunksReceived := []string{}

	if int(f.State) == lib.FileIncomplete {
		chunksNeeded = f.NumChunks - numChunks

		for _, c := range chunks {
			chunksReceived = append(chunksReceived, c.Hash)
		}
	}

	return &fileStatus{
		ID:             f.ID,
		State:          state,
		ChunksReceived: chunksReceived,
		ChunksNeeded:   chunksNeeded,
	}
}

func (f FileController) Status(e echo.Context) error {
	hash := e.Param("hash")

	file, err := models.Files(f.DB, Where("hash=$1", hash)).One()
	if err != nil {
		f.Debug("Failed to find file with hash", "hash", hash)
		return e.NoContent(http.StatusNotFound)
	}

	fs := statusForFile(f.DB, file)
	return e.JSON(http.StatusOK, fs)
}
