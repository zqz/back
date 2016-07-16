package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models/user"
)

func init() {
	lib.Connect()
}

func TestValidCredentials_Success(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		u := user.User{
			Username: "foo",
			Password: "bar",
		}

		u.Create(ex)

		exists := user.ValidCredentials(ex, "foo", "bar")
		a.Equal(exists, true)
	})
}

func TestValidCredentials_Failure(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		u := user.User{
			Username: "foo",
			Password: "bar",
		}

		u.Create(ex)

		exists := user.ValidCredentials(ex, "", "NOPE")
		a.Equal(exists, false)

		exists = user.ValidCredentials(ex, "foo", "")
		a.Equal(exists, false)

		exists = user.ValidCredentials(ex, "foo", "NOPE")
		a.Equal(exists, false)

		exists = user.ValidCredentials(ex, "NOPE", "bar")
		a.Equal(exists, false)
	})
}

func TestUsernameFree(t *testing.T) {
	t.Parallel()

	db.TxWrapper(func(ex db.Executor) {
		a := assert.New(t)

		// Below min length
		a.Equal(user.UsernameFree(ex, "a23"), false)
		// Over max length
		a.Equal(user.UsernameFree(ex, "abc3456789012345"), false)

		// Support unicode
		a.Equal(user.UsernameFree(ex, "世界"), false)
		a.Equal(user.UsernameFree(ex, "世界ada"), true)

		a.Equal(user.UsernameFree(ex, "foobar"), true)

		u := &user.User{
			FirstName: "Foo",
			LastName:  "Last",
			Username:  "foobar",
			Password:  "bar",
			Email:     "foo@bar.com",
		}

		u.Create(ex)

		// Username is taken
		a.Equal(user.UsernameFree(ex, "foobar"), false)
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
