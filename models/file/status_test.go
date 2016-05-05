package file

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/chunk"
)

func TestStatusForFile_Incomplete(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &File{Chunks: 2}
		f.Create(tx)

		s := StatusForFile(tx, f)

		a.Equal("Incomplete", s.State)
		a.Equal(f.ID, s.ID)
		a.Equal([]string{}, s.ChunksReceived)
		a.Equal(2, s.ChunksNeeded)
	})
}

func TestStatusForFile_IncompleteOneChunk(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &File{Chunks: 2}
		f.Create(tx)

		c := chunk.Chunk{FileID: f.ID, Position: 1, Hash: "c1"}
		c.Create(tx)

		s := StatusForFile(tx, f)

		a.Equal("Incomplete", s.State)
		a.Equal(f.ID, s.ID)
		a.Equal([]string{"c1"}, s.ChunksReceived)
		a.Equal(1, s.ChunksNeeded)
	})
}

func TestStatusForFile_IncompleteAllChunks(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &File{Chunks: 2}
		f.Create(tx)

		c := chunk.Chunk{FileID: f.ID, Position: 1, Hash: "c1"}
		c.Create(tx)
		c = chunk.Chunk{FileID: f.ID, Position: 0, Hash: "c2"}
		c.Create(tx)

		s := StatusForFile(tx, f)

		a.Equal("Incomplete", s.State)
		a.Equal(f.ID, s.ID)
		a.Equal([]string{"c2", "c1"}, s.ChunksReceived)
		a.Equal(0, s.ChunksNeeded)
	})
}
