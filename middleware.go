package main

// import (
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/zqzca/echo"
// )

// func CORSMiddleware() echo.MiddlewareFunc {
// 	return func(h echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c *echo.Context) error {
// 			req := c.Request()
// 			res := c.Response()

// 			res.Header().Add("Access-Control-Allow-Origin", "*")
// 			res.Header().Add(
// 				"Access-Control-Allow-Headers",
// 				"Access-Control-Request-Method, Content-Type, Authorization",
// 			)
// 			res.Header().Add(
// 				"Access-Control-Allow-Methods",
// 				"OPTIONS, POST, GET, PATCH, DELETE",
// 			)

// 			// I don't use OPTIONS so these only come from preflight requests.
// 			if req.Method == "OPTIONS" {
// 				return nil
// 			}

// 			if err := h(c); err != nil {
// 				c.Error(err)
// 			}

// 			return nil
// 		}
// 	}
// }

// // A JSON Web Token middleware
// func JWTAuth() echo.HandlerFunc {
// 	Bearer := "Bearer"
// 	Secret := os.Getenv("JWT_SECRET")

// 	return func(c *echo.Context) error {

// 		// Skip WebSocket
// 		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
// 			return nil
// 		}

// 		auth := c.Request().Header.Get("Authorization")
// 		l := len(Bearer)
// 		he := echo.NewHTTPError(http.StatusUnauthorized)

// 		if len(auth) > l+1 && auth[:l] == Bearer {
// 			t, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {

// 				// Always check the signing method
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 				}

// 				// Return the key for validation
// 				return []byte(Secret), nil
// 			})
// 			if err == nil && t.Valid {
// 				// Store token claims in echo.Context
// 				c.Set("claims", t.Claims)
// 				return nil
// 			}
// 		}
// 		return he
// 	}
// }
