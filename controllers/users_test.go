package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/models"
)

func TestUserCreateValid(t *testing.T) {
	models.TruncateUsers()

	a := assert.New(t)

	u := &models.User{
		Username:  "foo",
		Password:  "bar",
		Email:     "foo@bar.com",
		APIKey:    "123456",
		FirstName: "Foo",
		LastName:  "Bar",
	}

	res, c := post(u.String())

	UserCreate(c)

	a.Equal(201, res.Code)
	a.Equal(1, models.UserCount())

	u = models.LoadUser(res.Body)
	a.Equal("foo", u.Username)
	a.Equal("Foo", u.FirstName)
	a.Equal("Bar", u.LastName)
	a.NotEmpty(u.ID)
}

func TestUserGetValid(t *testing.T) {
	models.TruncateUsers()

	a := assert.New(t)

	u := CreateUser("foo", "bar")

	res, c := get(u.String())

	c.SetParam("id", u.ID)

	UserGet(c)

	a.Equal(200, res.Code)
	a.Equal(1, models.UserCount())

	u = models.LoadUser(res.Body)

	a.Equal("foo", u.Username)
}
