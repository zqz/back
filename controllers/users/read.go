package users

import (
	"net/http"

	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"
)

func (u Controller) Read(e echo.Context) error {
	id := e.Param("id")

	user, err := models.FindUser(u.DB, id)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, user)
}
