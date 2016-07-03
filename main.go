package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
	"github.com/zqzca/back/controllers/chunks"
	"github.com/zqzca/back/controllers/files"
	"github.com/zqzca/back/controllers/sessions"
	"github.com/zqzca/back/controllers/users"
	"github.com/zqzca/back/db"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
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

	connect()
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status} ${method} ${uri} - ${latency_human}, rx=${rx_bytes}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	// servers other static files
	// e.ServeDir("/assets", "assets")
	// e.ServeFile("/", "assets/index.html")
	// e.ServeFile("/favicon.ico", "assets/favicon.ico")
	e.Get("/d/:id", files.Download)

	// api := e.Group("/api")
	v1 := e.Group("/api/v1")

	// api.Get("/p2p/join/:id", controllers.P2PJoin)
	// api.Post("/p2p/join/:id", controllers.P2PJoinAnswer)
	// api.Get("/p2p/signaling", controllers.P2PWS)
	// Route
	// e.Get("/chunk/status", controllers.ChunkStatus)

	// Files
	v1.Get("/check/:hash", files.Status)
	v1.Get("/files", files.Index)
	v1.Get("/files/:file_id", files.Read)
	v1.Get("/files/:id/data", files.Download)
	v1.Post("/files", files.Create)
	v1.Post("/files/:id/process", files.Process)

	// Chunks
	v1.Post("/files/:file_id/chunks/:chunk_id", chunks.Write)
	v1.Get("/files/:file_id/chunks/:chunk_id", chunks.Read)

	// Users
	v1.Post("/users", users.Create)
	v1.Get("/username/valid", users.Valid)
	v1.Get("/users/:id", users.Read)

	// Sessions
	v1.Post("/sessions", sessions.Create)

	// r := api.Group("/users")
	// r.Use(JWTAuth())
	// r.Get("/:id", controllers.UserGet)
	// e.Patch("/users/:id", updateUser)
	// e.Delete("/users/:id", deleteUser)

	// e.ServeFile("/signin", "assets/signin.html")
	// e.ServeFile("/*", "assets/index.html")
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "assets",
		Browse: false,
		Index:  "index.html",
		HTML5:  true,
	}))

	// Start server
	e.Run(standard.New(":3001"))

}

func connect() error {
	open := os.Getenv("DATABASE_URL")

	if parsedURL, err := pq.ParseURL(open); err == nil && parsedURL != "" {
		open = parsedURL
	}

	con, err := sql.Open("postgres", open)

	if err != nil {
		fmt.Println(err)
	}

	db.Connection = con

	return err
}
