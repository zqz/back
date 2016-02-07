package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

type Session struct {
	Username string
	Password string
}

func (s Session) String() string {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(s)
	return buf.String()
}

type SessionError struct {
	Msg string `json:"error"`
}

func SessionCreate(c *echo.Context) error {
	s := &Session{}

	if err := c.Bind(s); err != nil {
		return err
	}

	if u, err := models.UserFindByLogin(s.Username, s.Password); err != nil {
		errors := &SessionError{err.Error()}
		return c.JSON(http.StatusUnauthorized, errors)
	} else {
		return c.JSON(http.StatusCreated, u)
	}
}
