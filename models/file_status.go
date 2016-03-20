package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FileStatus struct {
	ID             string   `json:"id"`
	State          string   `json:"state"`
	ChunksReceived []string `json:"chunks_received"`
	ChunksNeeded   int      `json:"chunks_needed,omitempty"`
}

func FileStatusForFile(f *File) *FileStatus {
	if f == nil {
		return nil
	}

	var state string

	switch f.State {
	case Incomplete:
		state = "Incomplete"
	case Assembling:
		state = "Assembling"
	case Processing:
		state = "Processing"
	case Finished:
		state = "Finished"
	}

	chunksNeeded := 0
	chunks, err := ChunksByFileID(f.ID)

	if err != nil {
		fmt.Println("fuck fuck")
		fmt.Println(err)
		return nil
	}

	numChunks := len(chunks)
	chunksReceived := make([]string, numChunks)

	for _, c := range chunks {
		chunksReceived = append(chunksReceived, c.Hash)
	}

	if f.State == Incomplete {
		chunksNeeded = f.Chunks - numChunks
	}

	return &FileStatus{
		ID:             f.ID,
		State:          state,
		ChunksReceived: chunksReceived,
		ChunksNeeded:   chunksNeeded,
	}
}

func (f *FileStatus) String() string {
	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(f)

	return buf.String()
}
