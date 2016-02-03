package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
)

func CreateUser(username string, password string) *models.User {
	u := &models.User{
		FirstName: "Tester",
		LastName:  "McTesterson",
		Email:     "foo@bar.com",
		APIKey:    "123456",
		Username:  username,
		Password:  password,
	}

	u.Save()

	return u
}

func CreateSessionRequest(username string, password string) string {
	s := Session{
		Username: username,
		Password: password,
	}

	return s.String()
}

func TestSessionCreateValid(t *testing.T) {
	models.TruncateUsers()

	a := assert.New(t)

	CreateUser("foo", "bar")

	res, c := post(CreateSessionRequest("foo", "bar"))

	SessionCreate(c)

	a.Equal(http.StatusCreated, res.Code)

	u := models.User{}
	json.NewDecoder(res.Body).Decode(&u)

	// Should return User struct
	a.Equal(u.Username, "foo")
	a.Empty(u.Password, "foo")
}

func TestSessionCreateInvalid(t *testing.T) {
	models.TruncateUsers()

	a := assert.New(t)

	res, c := post(CreateSessionRequest("foo", "bar"))

	SessionCreate(c)
	a.Equal(http.StatusUnauthorized, res.Code)

	s := SessionError{}
	json.NewDecoder(res.Body).Decode(&s)
	a.Equal("Invalid Credentials", s.Msg)
}
