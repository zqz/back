package chunk_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/chunk"
	"github.com/zqzca/back/models/file"
)

func createFile(tx *sql.Tx) *file.File {
	f := &file.File{}
	f.Create(tx)
	return f
}

func TestCreate(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := createFile(tx)

		c := &chunk.Chunk{
			Size:     123,
			Hash:     "foo",
			FileID:   f.ID,
			Position: 1,
		}

		c_err := c.Create(tx)

		// There should not be an error.
		a.Nil(c_err)

		// Postgres will assign an ID.
		a.NotEmpty(c.ID)
	})
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)
		f := createFile(tx)
		c := &chunk.Chunk{
			Size:     123,
			Hash:     "foo",
			FileID:   f.ID,
			Position: 1,
		}

		c.Create(tx)

		e, err := chunk.FindByID(tx, c.ID)

		// There should not be an error.
		a.Nil(err)

		a.Equal(e.ID, c.ID)
		a.Equal(e.Size, 123)
		a.Equal(e.Hash, "foo")
		a.Equal(e.FileID, f.ID)
		a.Equal(e.Position, 1)
		a.NotNil(e.CreatedAt)
		a.NotNil(e.UpdatedAt)
	})
}

func TestFindByFileID(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &file.File{}
		f.Create(tx)

		c1 := &chunk.Chunk{Hash: "a", FileID: f.ID, Position: 1}
		c1.Create(tx)
		c2 := &chunk.Chunk{Hash: "b", FileID: f.ID, Position: 1}
		c2.Create(tx)
		c3 := &chunk.Chunk{Hash: "c", FileID: f.ID, Position: 1}
		c3.Create(tx)

		chunks, err := chunk.FindByFileID(tx, f.ID)

		a.Nil(err)
		a.Equal(len(*chunks), 3)

		// Ordered
		a.Equal((*chunks)[0].Hash, "a")
		a.Equal((*chunks)[1].Hash, "b")
		a.Equal((*chunks)[2].Hash, "c")
	})
}

func TestHaveChunkForFile(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)
		f := &file.File{}
		f.Create(tx)

		haveChunk := chunk.HaveChunkForFile(tx, f.ID, 1)
		a.Equal(false, haveChunk)

		c1 := &chunk.Chunk{Hash: "a", FileID: f.ID, Position: 1}
		c1.Create(tx)

		haveChunk = chunk.HaveChunkForFile(tx, f.ID, 1)
		a.Equal(true, haveChunk)
	})
}
