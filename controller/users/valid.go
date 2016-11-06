package users

import (
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/zqzca/back/models"
	"github.com/labstack/echo"

	"github.com/vattle/sqlboiler/queries/qm"
)

// ValidateUsername checks to see if the username is in use
func (u Controller) ValidateUsername(e echo.Context) error {
	name := e.Param("name")

	count, err := models.Users(u.DB, qm.Where("username=?", name)).Count()
	if err != nil {
		log.Error("failed to get user from db", "err", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	if count > 0 {
		return e.NoContent(http.StatusNotAcceptable)
	}

	return e.NoContent(http.StatusOK)
}
