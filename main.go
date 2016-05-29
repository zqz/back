package main

import (
	"github.com/zqzca/back/controllers"

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

// func sshServer() {
// 	s := scp.NewSCPServer()
// 	s.ListenAndServe()
// }

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(CORSMiddleware())

	// servers other static files
	e.ServeDir("/assets", "assets")
	e.ServeFile("/", "assets/index.html")
	e.ServeFile("/favicon.ico", "assets/favicon.ico")
	e.Get("/d/:file_id", controllers.FileDownload)

	api := e.Group("/api")

	api.Get("/p2p/join/:id", controllers.P2PJoin)
	api.Post("/p2p/join/:id", controllers.P2PJoinAnswer)
	api.Get("/p2p/signaling", controllers.P2PWS)
	// Route
	// e.Get("/chunk/status", controllers.ChunkStatus)
	api.Get("/file/:hash", controllers.FileStatus)
	api.Get("/files", controllers.FileIndex)
	api.Get("/files/status/:hash", controllers.FileStatus)
	api.Get("/files/download/:file_id", controllers.FileDownload)
	api.Post("/sessions", controllers.SessionCreate)
	api.Post("/files", controllers.FileCreate)
	api.Post("/chunks", controllers.ChunkCreate)
	api.Post("/users", controllers.UserCreate)
	api.Get("/users/validateUsername", controllers.UserNameValid)

	r := api.Group("/users")
	r.Use(JWTAuth())
	r.Get("/:id", controllers.UserGet)
	// e.Patch("/users/:id", updateUser)
	// e.Delete("/users/:id", deleteUser)

	e.ServeFile("/signin", "assets/signin.html")
	e.ServeFile("/*", "assets/index.html")

	// Start server
	e.Run(":3001")
}
