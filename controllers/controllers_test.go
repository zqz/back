package controllers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/zqzca/back/models"
	"github.com/labstack/echo"
)

func request(method string, path string, jsonRequest string) (*httptest.ResponseRecorder, *echo.Context) {
	params := strings.NewReader(jsonRequest)
	req, _ := http.NewRequest(echo.POST, "/", params)
	req.Header.Add("Content-Type", "application/json")

	e := echo.New()
	rec := httptest.NewRecorder()
	res := echo.NewResponse(rec, e)

	c := echo.NewContext(req, res, e)

	return rec, c
}

func get(r string) (*httptest.ResponseRecorder, *echo.Context) {
	return request("GET", "/", r)
}

func post(r string) (*httptest.ResponseRecorder, *echo.Context) {
	return request("POST", "/", r)
}

func CreateUser(tx *sql.Tx, username string, password string) *models.User {
	u := &models.User{
		FirstName: "Tester",
		LastName:  "McTesterson",
		Email:     "foo@bar.com",
		Username:  username,
		Password:  password,
	}

	u.Create(tx)

	return u
}
