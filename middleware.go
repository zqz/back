package main

import "github.com/labstack/echo"

func CORSMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			res := c.Response()

			res.Header().Add("Access-Control-Allow-Origin", "*")
			res.Header().Add(
				"Access-Control-Allow-Headers",
				"Access-Control-Request-Method, Content-Type, Authorization",
			)
			res.Header().Add(
				"Access-Control-Allow-Methods",
				"OPTIONS, POST, GET, PATCH, DELETE",
			)

			// I don't use OPTIONS so these only come from preflight requests.
			if req.Method == "OPTIONS" {
				return nil
			}

			if err := h(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
