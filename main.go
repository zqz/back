package main

import (
	"net/http"
	"strconv"

	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	user struct {
		ID   int
		Name string
	}
)

var (
	users = map[int]*user{}
	seq   = 1
)

//----------
// Handlers
//----------

func createUser(c *echo.Context) error {
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func updateUser(c *echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c *echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	db := db.DatabaseConnect()

	models.SetDB(db)
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(CORSMiddleware())

	// Route
	e.Post("/sessions", controllers.SessionCreate)
	e.Post("/users", controllers.UserCreate)
	e.Get("/users/:id", controllers.UserGet)
	e.Patch("/users/:id", updateUser)
	e.Delete("/users/:id", deleteUser)

	// Start server
	e.Run(":3001")
}
