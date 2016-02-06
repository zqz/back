package controllers

import (
	"fmt"
	"net/http"

	"github.com/DylanJ/echo"
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

func UserGet(c *echo.Context) error {
	id := c.Param("id")
	if u, err := models.UserFind(id); err != nil {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.JSON(http.StatusOK, u)
	}
}

func UserIndex(c *echo.Context) error {
	return nil
}
