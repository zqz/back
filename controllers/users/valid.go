package users

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/user"
)

type UserError struct {
	Msg string `json:"error"`
}

func Valid(c echo.Context) error {
	tx := db.StartTransaction()
	defer tx.Rollback()
	name := c.Param("name")

	if user.UsernameFree(tx, name) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusNotAcceptable)
	}
}
