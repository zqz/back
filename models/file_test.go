package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Exists(t *testing.T) {
	truncate("files")

	a := assert.New(t)

	file, _ := FileFindByHash("foo")

	a.Nil(file)

	f := &File{
		Size:   123,
		Hash:   "foo",
		Done:   false,
		Chunks: 1,
		Type:   "image/jpg",
	}

	f.Save()

	file, _ = FileFindByHash("foo")

	a.NotNil(file)
}
