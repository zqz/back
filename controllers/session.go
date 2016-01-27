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

	if u, _ := models.UserFindByLogin(s.Username, s.Password); u != nil {
		return c.NoContent(http.StatusUnauthorized)
	} else {
		return c.JSON(http.StatusCreated, u)
	}
}
