package users

import (
	"fmt"
	"net/http"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/user"
	"github.com/zqzca/echo"
)

func Create(c echo.Context) error {
	u := &user.User{}

	if err := c.Bind(u); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if u.Valid() {
		u.Create(db.Connection)
		return c.JSON(http.StatusCreated, u)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}
