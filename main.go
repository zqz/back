package main

import (
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

func main() {
	db := db.DatabaseConnect()

	models.SetDB(db)
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(CORSMiddleware())

	// servers other static files
	e.ServeDir("/assets", "assets")
	e.ServeFile("/", "assets/index.html")
	e.ServeFile("/favicon.ico", "assets/favicon.ico")

	// Route
	// e.Get("/chunk/status", controllers.ChunkStatus)
	e.Get("/d/:file_id", controllers.FileDownload)
	e.Get("/files", controllers.FileIndex)
	e.Get("/files/status/:hash", controllers.FileStatus)
	e.Get("/files/download/:file_id", controllers.FileDownload)
	e.Post("/sessions", controllers.SessionCreate)
	e.Post("/files", controllers.FileCreate)
	e.Post("/chunks", controllers.ChunkCreate)
	e.Post("/users", controllers.UserCreate)
	e.Get("/users/validateUsername", controllers.UserNameValid)

	r := e.Group("/users")
	r.Use(JWTAuth())
	r.Get("/:id", controllers.UserGet)
	// e.Patch("/users/:id", updateUser)
	// e.Delete("/users/:id", deleteUser)

	// Start server
	e.Run(":3001")
}
