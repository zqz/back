package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/file"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		err := f.Create(ex)

		// There should not be an error.
		a.Nil(err)

		// Postgres will assign an ID.
		a.NotEmpty(f.ID)
	})
}

func TestFindByHash(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			Hash:  "foobar",
			State: file.Incomplete,
		}

		f.Create(ex)

		e, err := file.FindByHash(ex, f.Hash)

		// There should not be an error.
		a.Nil(err)

		a.Equal(f.ID, e.ID)
	})
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		f.Create(ex)

		e, err := file.FindByID(ex, f.ID)

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
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		f.Create(ex)
		a.Equal(f.State, file.Incomplete)

		f.SetState(ex, file.Processing)
		a.Equal(f.State, file.Processing)
	})
}

func TestPagination(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		f := &file.File{
			Size:  1,
			Name:  "doo",
			State: file.Incomplete,
		}

		f.Create(ex)

		files, err := file.Pagination(ex, 0, 10)
		a.Nil(err)
		a.NotNil((*files)[0])
		a.Equal((*files)[0].ID, f.ID)

		f.SetState(ex, file.Processing)
		a.Equal(f.State, file.Processing)
	})
}
