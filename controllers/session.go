package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models/user"
)

type session struct {
	Username string
	Password string
}

type sessionError struct {
	Msg string `json:"error"`
}

func SessionCreate(c *echo.Context) error {
	tx := StartTransaction()
	defer tx.Rollback()
	s := &session{}

	if err := c.Bind(s); err != nil {
		return err
	}

	if user.ValidCredentials(tx, s.Username, s.Password) {
		u, _ := user.FindByUsername(tx, s.Username)
		return c.JSON(http.StatusCreated, u)
	} else {
		errors := &sessionError{"Invalid Credentials"}
		return c.JSON(http.StatusUnauthorized, errors)
	}
}
