package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/models/user"
)

func CreateSessionRequest(username string, password string) string {
	return fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
}

func TestSessionCreateValid(t *testing.T) {
	db, err := models.Connection()
	a := assert.New(t)

	u := &user.User{
		Username: "foo",
		Password: "bar",
	}

	err = u.Create(db)

	a.NoError(err)

	res, c := post(CreateSessionRequest("foo", "bar"))

	SessionCreate(c)
	u.Delete(db)

	a.Equal(http.StatusCreated, res.Code)

	u = &user.User{}
	json.NewDecoder(res.Body).Decode(u)

	// Should return User struct
	a.Equal(u.Username, "foo")
	a.Empty(u.Password)
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
