package file

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zqzca/back/models/chunk"
)

type Status struct {
	ID             string   `json:"id"`
	State          string   `json:"state"`
	ChunksReceived []string `json:"chunks_received"`
	ChunksNeeded   int      `json:"chunks_needed,omitempty"`
}

func StatusForFile(tx *sql.Tx, f *File) *Status {
	if f == nil {
		return nil
	}

	var state string

	switch f.State {
	case Incomplete:
		state = "Incomplete"
	case Processing:
		state = "Processing"
	case Finished:
		state = "Finished"
	}

	chunksNeeded := 0
	chunks, err := chunk.FindByFileID(tx, f.ID)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	numChunks := len(*chunks)
	chunksReceived := []string{}

	for _, c := range *chunks {
		chunksReceived = append(chunksReceived, c.Hash)
	}

	if f.State == Incomplete {
		chunksNeeded = f.Chunks - numChunks
	}

	return &Status{
		ID:             f.ID,
		State:          state,
		ChunksReceived: chunksReceived,
		ChunksNeeded:   chunksNeeded,
	}
}

func (f *Status) String() string {
	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(f)

	return buf.String()
}
