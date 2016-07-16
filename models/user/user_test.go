package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/user"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		u := &user.User{
			Username: "foo",
			Password: "bar",
		}

		err := u.Create(ex)

		a.Nil(err)
		a.NotNil(u.ID)
	})
}

func TestFindByUsername_Success(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)
		u := user.User{
			Username: "foo",
			Password: "bar",
		}
		u.Create(ex)

		e, err := user.FindByUsername(ex, "foo")

		a.Nil(err)
		a.NotNil(e)
		a.Equal(e.ID, u.ID)
		a.Equal(e.Username, u.Username)
		a.Equal(e.FirstName, u.FirstName)
		a.Equal(e.LastName, u.LastName)
		a.Equal(e.Phone, u.Phone)
		a.Equal(e.Email, u.Email)
		a.Equal(e.Banned, u.Banned)
		a.NotNil(e.CreatedAt)
		a.NotNil(e.UpdatedAt)
	})
}

func TestFindByUsername_Failure(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)
		u := user.User{
			Username: "foo",
			Password: "bar",
		}
		u.Create(ex)

		e, err := user.FindByUsername(ex, "bar")

		a.NotNil(err)
		a.Empty(e.ID)
	})
}

func TestSetPassword(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)
		u := user.User{
			Username: "foo",
			Password: "first",
		}
		u.Create(ex)

		exists := user.ValidCredentials(ex, "foo", "first")
		a.Equal(exists, true)

		u.SetPassword(ex, "second")
		exists = user.ValidCredentials(ex, "foo", "first")
		a.Equal(exists, false)

		exists = user.ValidCredentials(ex, "foo", "second")
		a.Equal(exists, true)
	})
}
