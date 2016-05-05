package user_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/user"
)

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
