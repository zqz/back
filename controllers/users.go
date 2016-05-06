package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models/user"
)

type UserError struct {
	Msg string `json:"error"`
}

func UserCreate(c *echo.Context) error {
	u := &user.User{}

	if err := c.Bind(u); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if u.Valid() {
		tx := StartTransaction()
		u.Create(tx)
		tx.Commit()
		return c.JSON(http.StatusCreated, u)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}

func UserGet(c *echo.Context) error {
	// tx := StartTransaction()
	// defer tx.Rollback()
	// id := GetParam(c, "id")

	// if u, err := user.FindByID(tx, id); err != nil {
	// 	errors := &UserError{err.Error()}
	// 	return c.JSON(http.StatusOK, u)
	// 	return c.JSON(http.StatusNotFound, errors)
	// } else {
	// 	return c.JSON(http.StatusOK, u)
	// }
	return nil
}

func UserIndex(c *echo.Context) error {
	return nil
}

func UserNameValid(c *echo.Context) error {
	tx := StartTransaction()
	defer tx.Rollback()
	name := GetParam(c, "name")

	if user.UsernameFree(tx, name) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusNotAcceptable)
	}
}
