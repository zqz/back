package user_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/user"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		u := &user.User{
			Username: "foo",
			Password: "bar",
		}

		err := u.Create(tx)

		a.Nil(err)
		a.NotNil(u.ID)
	})
}

func TestFindByUsername_Success(t *testing.T) {
	t.Parallel()

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)
		u := user.User{
			Username: "foo",
			Password: "bar",
		}
		u.Create(tx)

		e, err := user.FindByUsername(tx, "foo")

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

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)
		u := user.User{
			Username: "foo",
			Password: "bar",
		}
		u.Create(tx)

		e, err := user.FindByUsername(tx, "bar")

		a.NotNil(err)
		a.Empty(e.ID)
	})
}

func TestSetPassword(t *testing.T) {
	t.Parallel()

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)
		u := user.User{
			Username: "foo",
			Password: "first",
		}
		u.Create(tx)

		exists := user.ValidCredentials(tx, "foo", "first")
		a.Equal(exists, true)

		u.SetPassword(tx, "second")
		exists = user.ValidCredentials(tx, "foo", "first")
		a.Equal(exists, false)

		exists = user.ValidCredentials(tx, "foo", "second")
		a.Equal(exists, true)
	})
}
