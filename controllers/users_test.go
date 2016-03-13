package controllers_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	. "github.com/zqzca/back/controllers"
	"github.com/zqzca/back/models"
)

func TestUserCreateValid(t *testing.T) {
	models.Truncate("users")

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
	models.Truncate("users")

	a := assert.New(t)

	u := CreateUser("foo", "bar")

	res, c := get("users/id")

	GetParam = func(c *echo.Context, key string) string {
		return u.ID
	}

	UserGet(c)

	a.Equal(200, res.Code)
	a.Equal(1, models.UserCount())

	u = models.LoadUser(res.Body)

	a.Equal("foo", u.Username)
}

func TestUserNameValid(t *testing.T) {
	models.Truncate("Users")

	a := assert.New(t)

	res, c := get("")

	// Stub params to eq username
	GetParam = func(c *echo.Context, key string) string {
		return "foobar"
	}

	UserNameValid(c)
	a.Equal(200, res.Code)

	u := CreateUser("foobar", "bar")
	u.Save()

	res, c = get("")

	UserNameValid(c)
	a.Equal(406, res.Code)
}
