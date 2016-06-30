package users

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/user"
)

func Create(c *echo.Context) error {
	u := &user.User{}

	if err := c.Bind(u); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if u.Valid() {
		tx := db.StartTransaction()
		u.Create(tx)
		tx.Commit()
		return c.JSON(http.StatusCreated, u)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}
