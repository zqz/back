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

func TestValidCredentials_Success(t *testing.T) {
	t.Parallel()

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		u := user.User{
			Username: "foo",
			Password: "bar",
		}

		u.Create(tx)

		exists := user.ValidCredentials(tx, "foo", "bar")
		a.Equal(exists, true)
	})
}

func TestValidCredentials_Failure(t *testing.T) {
	t.Parallel()

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		u := user.User{
			Username: "foo",
			Password: "bar",
		}

		u.Create(tx)

		exists := user.ValidCredentials(tx, "", "NOPE")
		a.Equal(exists, false)

		exists = user.ValidCredentials(tx, "foo", "")
		a.Equal(exists, false)

		exists = user.ValidCredentials(tx, "foo", "NOPE")
		a.Equal(exists, false)

		exists = user.ValidCredentials(tx, "NOPE", "bar")
		a.Equal(exists, false)
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

func TestUsernameFree(t *testing.T) {
	t.Parallel()

	models.TxWrapper(func(tx *sql.Tx) {
		a := assert.New(t)

		// Below min length
		a.Equal(user.UsernameFree(tx, "a23"), false)
		// Over max length
		a.Equal(user.UsernameFree(tx, "abc3456789012345"), false)

		// Support unicode
		a.Equal(user.UsernameFree(tx, "世界"), false)
		a.Equal(user.UsernameFree(tx, "世界ada"), true)

		a.Equal(user.UsernameFree(tx, "foobar"), true)

		u := &user.User{
			FirstName: "Foo",
			LastName:  "Last",
			Username:  "foobar",
			Password:  "bar",
			Email:     "foo@bar.com",
		}

		u.Create(tx)

		// Username is taken
		a.Equal(user.UsernameFree(tx, "foobar"), false)
	})
}

func TestValid(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	u := &user.User{}

	a.Equal(false, u.Valid(), "Empty user should not be valid.")
	u = &user.User{
		FirstName: "John",
		LastName:  "Carmack",
		Username:  "johnc",
		Email:     "jc@idsoftware.com",
	}

	a.Equal(true, u.Valid(), "Should be a valid User.")

	u.Email = ""
	a.Equal(false, u.Valid(), "Missing email should be invalid.")
	u.Email = "jc@idsoftware.com"

	u.FirstName = ""
	a.Equal(false, u.Valid(), "Missing first name should be invalid.")
	u.FirstName = "John"

	u.LastName = ""
	a.Equal(false, u.Valid(), "Missing last name should be invalid.")
	u.LastName = "Carmack"

	u.Username = ""
	a.Equal(false, u.Valid(), "Missing username should be invalid.")
	u.Username = "johnc"
}
