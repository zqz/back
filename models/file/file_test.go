package file_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/file"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		err := f.Create(tx)

		// There should not be an error.
		a.Nil(err)

		// Postgres will assign an ID.
		a.NotEmpty(f.ID)
	})
}

func TestFindByHash(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			Hash:  "foobar",
			State: file.Incomplete,
		}

		f.Create(tx)

		e, err := file.FindByHash(tx, f.Hash)

		// There should not be an error.
		a.Nil(err)

		a.Equal(f.ID, e.ID)
	})
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		f.Create(tx)

		e, err := file.FindByID(tx, f.ID)

		// There should not be an error.
		a.Nil(err)

		a.Equal(e.ID, f.ID)
		a.Equal(e.Size, f.Size)
		a.Equal(e.Chunks, f.Chunks)
		a.Equal(e.Name, f.Name)
		a.Equal(e.Type, f.Type)
		a.Equal(e.State, f.State)
		a.NotNil(e.CreatedAt)
		a.NotNil(e.UpdatedAt)
	})
}

func TestSetState(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		f.Create(tx)
		a.Equal(f.State, file.Incomplete)

		f.SetState(tx, file.Processing)
		a.Equal(f.State, file.Processing)
	})
}
