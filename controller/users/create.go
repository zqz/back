package users

import (
	"fmt"
	"net/http"

	"github.com/zqzca/back/models"
	"github.com/labstack/echo"
)

// Create creates a new user
func (u Controller) Create(c echo.Context) error {
	user := &models.User{}

	if err := c.Bind(u); err != nil {
		fmt.Println(err.Error())
		return err
	}

	//if !u.Valid() {
	//	return c.NoContent(http.StatusBadRequest)
	//}

	if err := user.Insert(u.DB); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, u)
}
