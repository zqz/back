package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/zqzca/back/models"
)

func TestFile_Create(t *testing.T) {
	Truncate("files")

	a := assert.New(t)

	f := &File{
		Name:   "foo",
		Size:   123,
		Hash:   "foo",
		Chunks: 1,
		Type:   "image/jpg",
	}

	f.Save()

	a.NotEmpty(f.ID)
}

func TestFile_FindByHash(t *testing.T) {
	Truncate("files")

	a := assert.New(t)

	file, _ := FileFindByHash("foo")

	a.Nil(file)

	f := &File{
		Name:   "foo",
		Size:   123,
		Hash:   "foo",
		Chunks: 1,
		Type:   "image/jpg",
	}

	f.Save()

	file, _ = FileFindByHash("foo")

	a.NotNil(file)
}

func TestFile_Status(t *testing.T) {

}
