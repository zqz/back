package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

func SessionCreate(c *echo.Context) error {
	s := &models.Session{}
	if err := c.Bind(s); err != nil {
		fmt.Println(err)
		return err
	}

	if s.Username == "foo" && s.Password == "bar" {
		return c.JSON(http.StatusCreated, s)
	} else {
		return c.NoContent(http.StatusUnauthorized)
	}
}
