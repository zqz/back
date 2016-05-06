package controllers_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/zqzca/back/models"
)

func CreateSessionRequest(username string, password string) string {
	return fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
}

func TestSessionCreateValid(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		// a := assert.New(t)

		// CreateUser(tx, "foo", "bar")

		// res, c := post(CreateSessionRequest("foo", "bar"))

		// SessionCreate(c)

		// a.Equal(http.StatusCreated, res.Code)

		// u := user.User{}
		// json.NewDecoder(res.Body).Decode(&u)

		// // Should return User struct
		// a.Equal(u.Username, "foo")
		// a.Empty(u.Password, "foo")
	})
}

func TestSessionCreateInvalid(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		// a := assert.New(t)

		// res, c := post(CreateSessionRequest("foo", "bar"))

		// SessionCreate(c)
		// a.Equal(http.StatusUnauthorized, res.Code)

		// // s := SessionError{}
		// // json.NewDecoder(res.Body).Decode(&s)
		// a.Equal("Invalid Credentials", res.Body)
	})
}
