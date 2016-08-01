package sessions

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"

	. "github.com/nullbio/sqlboiler/boil/qm"
)

const (
	bcryptCost = bcrypt.MinCost
)

var (
	errInvalidCredentials = []byte(`{"err":"invalid credentials"}`)
)

type userSession struct {
	Username string
	Password string
}

func (s SessionsController) Create(e echo.Context) error {
	session := &userSession{}

	if err := e.Bind(s); err != nil {
		return err
	}

	user, err := models.Users(s.DB, Select("hash"), Where("username=$1", session.Username)).One()
	if err != nil {
		s.Error("failed to fetch user", "err", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(session.Password)); err != nil {
		return e.JSONBlob(http.StatusUnauthorized, errInvalidCredentials)
	}

	return e.NoContent(http.StatusOK)
}
