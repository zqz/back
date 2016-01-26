package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
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

func TestSessionCreateValid(t *testing.T) {
	a := assert.New(t)
	res, c := request(
		"GET", "/",
		`{"username": "foo", "password": "bar"}`,
	)

	SessionCreate(c)

	// Should be Unauthorized
	a.Equal(201, res.Code, "Should be successful")
	a.Equal(
		`{"username": "foo"}`,
		res.Body.String(),
		"Should return user json",
	)
}
