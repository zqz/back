package thumbnail_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/file"
	"github.com/zqzca/back/models/thumbnail"
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
		t := &thumbnail.Thumbnail{
			Size:   123,
			Hash:   "foo",
			FileID: f.ID,
		}

		err := t.Create(ex)

		// There should not be an error.
		a.Nil(err)

		// Postgres will assign an ID.
		a.NotEmpty(t.ID)
	})
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)
		f := createFile(ex)
		t := &thumbnail.Thumbnail{
			Size:   123,
			Hash:   "foo",
			FileID: f.ID,
		}

		t.Create(ex)

		e, err := thumbnail.FindByID(ex, t.ID)

		a.Nil(err)

		a.Equal(e.ID, t.ID)
		a.Equal(e.Size, 123)
		a.Equal(e.Hash, "foo")
		a.NotNil(e.CreatedAt)
		a.NotNil(e.UpdatedAt)
	})
}
