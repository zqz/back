package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

type UserError struct {
	msg string
}

func UserCreate(c *echo.Context) error {
	u := &models.User{}

	if err := c.Bind(u); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if u.Valid() {
		u.Save()
		return c.JSON(http.StatusCreated, u)
	} else {
		return c.JSON(http.StatusBadRequest, u.Errors())
	}
}

func UserIndex(c *echo.Context) error {

	return nil

}
