package thumbnail_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/file"
	"github.com/zqzca/back/models/thumbnail"
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
		t := &thumbnail.Thumbnail{
			Size:   123,
			Hash:   "foo",
			FileID: f.ID,
		}

		err := t.Create(tx)

		// There should not be an error.
		a.Nil(err)

		// Postgres will assign an ID.
		a.NotEmpty(t.ID)
	})
}

func TestGetByFileID(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)
		f := createFile(tx)
		t := &thumbnail.Thumbnail{
			Size:   123,
			Hash:   "foo",
			FileID: f.ID,
		}

		t.Create(tx)

		e, err := thumbnail.FindByFileID(tx, t.ID)

		a.Nil(err)

		a.Equal(e.ID, t.ID)
		a.Equal(e.Size, 123)
		a.Equal(e.Hash, "foo")
		a.NotNil(e.CreatedAt)
		a.NotNil(e.UpdatedAt)
	})
}
