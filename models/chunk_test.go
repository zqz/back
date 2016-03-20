package models_test

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	. "github.com/zqzca/back/models"
)

func TestChunk_CreateWithFile(t *testing.T) {
	Truncate("chunks")

	a := assert.New(t)

	f := createFile()

	c := &Chunk{
		Size:     123,
		Hash:     "foo",
		FileID:   f.ID,
		Position: 1,
	}

	a.Equal(true, c.Save())

	a.NotEmpty(c.ID)

	Truncate("files")
}

func TestChunk_CreateWithoutFile(t *testing.T) {
	Truncate("chunks")

	a := assert.New(t)

	c := &Chunk{
		Size:     123,
		Hash:     "foo",
		FileID:   uuid.NewV4().String(),
		Position: 1,
	}

	a.Equal(false, c.Save())

	a.Empty(c.ID)
}
