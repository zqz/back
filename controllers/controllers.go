package controllers

import "github.com/labstack/echo"

func Param(c *echo.Context, key string) string {
	return c.Param(key)
}

var GetParam = Param
