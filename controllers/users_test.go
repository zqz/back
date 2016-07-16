package controllers_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/controllers/users"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/user"
)

func TestUserCreateValid(t *testing.T) {
	t.Parallel()
	db.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		u := &user.User{
			Username:  "foo",
			Password:  "bar",
			Email:     "foo@bar.com",
			FirstName: "Foo",
			LastName:  "Bar",
		}

		res, c := post(lib.ToJSON(u))

		users.Create(c)

		a.Equal(201, res.Code)

		u = models.LoadUser(res.Body)
		a.Equal("foo", u.Username)
		a.Equal("Foo", u.FirstName)
		a.Equal("Bar", u.LastName)
		a.NotEmpty(u.ID)
	})
}

func TestUserGetValid(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		// a := assert.New(t)

		// u := CreateUser(tx, "foo", "bar")

		// res, c := get("users/id")

		// GetParam = func(c *echo.Context, key string) string {
		// 	return u.ID
		// }

		// UserGet(c)

		// a.Equal(200, res.Code)
		// a.Equal(1, models.UserCount())

		// u = models.LoadUser(res.Body)

		// a.Equal("foo", u.Username)
	})
}

func TestUserNameValid(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {

		// 		a := assert.New(t)

		// 		res, c := get("")

		// 		// Stub params to eq username
		// 		GetParam = func(c *echo.Context, key string) string {
		// 			return "foobar"
		// 		}

		// 		UserNameValid(c)
		// 		a.Equal(200, res.Code)

		// 		u := CreateUser(tx, "foobar", "bar")
		// 		u.Create(tx)

		// 		res, c = get("")

		// 		UserNameValid(c)
		// 		a.Equal(406, res.Code)
	})
}
