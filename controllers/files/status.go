package files

import (
	"fmt"
	"net/http"

	"github.com/zqzca/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/chunk"
	"github.com/zqzca/back/models/file"
)

type fileStatus struct {
	ID             string   `json:"id"`
	State          string   `json:"state"`
	ChunksReceived []string `json:"chunks_received,omitempty"`
	ChunksNeeded   int      `json:"chunks_needed,omitempty"`
}

func statusForFile(ex db.Executor, f *file.File) *fileStatus {
	var state string

	switch f.State {
	case file.Incomplete:
		state = "incomplete"
	case file.Processing:
		state = "processing"
	case file.Finished:
		state = "finished"
	}

	chunksNeeded := 0
	chunks, err := chunk.FindByFileID(ex, f.ID)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	numChunks := len(*chunks)
	chunksReceived := []string{}

	if f.State == file.Incomplete {
		chunksNeeded = f.Chunks - numChunks

		for _, c := range *chunks {
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

func Status(c echo.Context) error {
	hash := c.Param("hash")
	tx := db.Connection

	f, err := file.FindByHash(tx, hash)

	if err != nil || f == nil {
		fmt.Println("Failed to find file with hash:", hash)
		return c.NoContent(http.StatusNotFound)
	}

	fs := statusForFile(tx, f)
	return c.JSON(http.StatusOK, fs)
}
