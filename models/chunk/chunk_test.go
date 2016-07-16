package chunk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/chunk"
	"github.com/zqzca/back/models/file"
)

func init() {
	lib.Connect()
}

func createFile(ex db.Executor) *file.File {
	f := &file.File{}
	f.Create(ex)
	return f
}

func TestCreate(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := createFile(ex)

		c := &chunk.Chunk{
			Size:     123,
			Hash:     "foo",
			FileID:   f.ID,
			Position: 1,
		}

		c_err := c.Create(ex)

		// There should not be an error.
		a.Nil(c_err)

		// Postgres will assign an ID.
		a.NotEmpty(c.ID)
	})
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)
		f := createFile(ex)
		c := &chunk.Chunk{
			Size:     123,
			Hash:     "foo",
			FileID:   f.ID,
			Position: 1,
		}

		c.Create(ex)

		e, err := chunk.FindByID(ex, c.ID)

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
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := &file.File{}
		f.Create(ex)

		c1 := &chunk.Chunk{Hash: "a", FileID: f.ID, Position: 1}
		c1.Create(ex)
		c2 := &chunk.Chunk{Hash: "b", FileID: f.ID, Position: 1}
		c2.Create(ex)
		c3 := &chunk.Chunk{Hash: "c", FileID: f.ID, Position: 1}
		c3.Create(ex)

		chunks, err := chunk.FindByFileID(ex, f.ID)

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
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)
		f := &file.File{}
		f.Create(ex)

		haveChunk := chunk.HaveChunkForFile(ex, f.ID, 1)
		a.Equal(false, haveChunk)

		c1 := &chunk.Chunk{Hash: "a", FileID: f.ID, Position: 1}
		c1.Create(ex)

		haveChunk = chunk.HaveChunkForFile(ex, f.ID, 1)
		a.Equal(true, haveChunk)
	})
}
