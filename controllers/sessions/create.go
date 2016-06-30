package sessions

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/user"
)

type session struct {
	Username string
	Password string
}

type sessionError struct {
	Msg string `json:"error"`
}

func Create(c *echo.Context) error {
	s := &session{}

	if err := c.Bind(s); err != nil {
		return err
	}

	if user.ValidCredentials(db.Connection, s.Username, s.Password) {
		u, _ := user.FindByUsername(db.Connection, s.Username)
		return c.JSON(http.StatusCreated, u)
	} else {
		errors := &sessionError{"Invalid Credentials"}
		return c.JSON(http.StatusUnauthorized, errors)
	}
}
