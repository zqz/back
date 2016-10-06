package app

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/controllers/chunks"
	"github.com/zqzca/back/controllers/dashboard"
	"github.com/zqzca/back/controllers/files"
	"github.com/zqzca/back/controllers/p2p"
	"github.com/zqzca/back/controllers/sessions"
	"github.com/zqzca/back/controllers/thumbnails"
	"github.com/zqzca/back/controllers/users"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/scp"
	"github.com/zqzca/back/ws"
	"github.com/zqzca/echo"

	"github.com/rsc/letsencrypt"
	"github.com/zqzca/echo/engine"
	"github.com/zqzca/echo/engine/standard"
	"github.com/zqzca/echo/middleware"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func redirect() {
	http.ListenAndServe(":3001", http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://"+req.Host+req.RequestURI, http.StatusMovedPermanently)
		},
	))
}

var config Config

// Run the application, start http and scp server.
func Run(appConfig Config) {
	config = appConfig

	// Connect to DB
	db, err := lib.Connect()

	if err != nil {
		fmt.Printf("Failed to connect to db")
		return
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status} ${method} ${uri} - ${latency_human}, rx=${rx_bytes}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Logging
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{}

	// Websockets
	ws := ws.NewServer()

	// Shared dependencies between all controllers
	deps := controllers.Dependencies{
		Fs:     afero.NewOsFs(),
		Logger: log,
		DB:     db,
		WS:     ws,
	}

	ws.Dependencies = &deps
	ws.Start()

	e.Get("/", Index)

	// Base API Group
	v1 := e.Group("/api/v1")

	// Files
	files := files.Controller{Dependencies: deps}
	e.Get("/d/:slug", files.Download)
	v1.Get("/check/:hash", files.Status)
	v1.Get("/files", files.Index)
	v1.Get("/files/:slug", files.Read)
	v1.Get("/files/:slug/data", files.Download)
	v1.Delete("/files/:slug/delete", files.Delete)
	v1.Post("/files", files.Create)
	v1.Get("/files/:slug/process", files.Process)

	// Thumbnail
	thumbnails := thumbnails.Controller{Dependencies: deps}
	v1.Get("/thumbnails/:id", thumbnails.Download)

	// Chunks
	chunks := chunks.NewController(deps)
	v1.Post("/files/:file_id/chunks/:chunk_id/:hash/:ws_id", chunks.Write)

	// Users
	users := users.Controller{Dependencies: deps}
	v1.Post("/users", users.Create)
	v1.Get("/username/valid", users.ValidateUsername)
	v1.Get("/users/:id", users.Read)

	// Sessions
	sessions := sessions.Controller{Dependencies: deps}
	v1.Post("/sessions", sessions.Create)

	// P2P
	v1.Get("/p2p/signaling", standard.WrapHandler(http.HandlerFunc(p2p.Signaling())))
	v1.Get("/p2p/:id", p2p.Join)
	v1.Post("/p2p/:id", p2p.Answer)

	// Dashboard
	dash := dashboard.Controller{Dependencies: deps}
	v1.Get("/dashboard", dash.Index)

	e.Static("/assets", "assets")
	e.Get("/*", Index)

	// WebSockets
	e.Get("/ws", standard.WrapHandler(ws.Endpoint()))

	// Horribleness for optional letsencrypt stuff
	var s *standard.Server

	if config.Secure {
		fmt.Println("Running in Secure Mode")
		var m letsencrypt.Manager
		if err := m.CacheFile("certs/letsencrypt.cache"); err != nil {
			log.Fatal(err)
		}

		cfg := &tls.Config{
			GetCertificate: m.GetCertificate,
		}

		engineConfig := engine.Config{
			Address:   config.HTTPBindAddr,
			TLSConfig: cfg,
		}

		s = standard.WithConfig(engineConfig)
		http2.ConfigureServer(s.Server, &http2.Server{})
		deps.Info("Listening for HTTP2 connections", "addr", config.HTTPBindAddr)
		go redirect()
	} else {
		deps.Info("Listening for HTTP connections", "addr", config.HTTPBindAddr)
		s = standard.New(config.HTTPBindAddr)
	}

	// Start SCP
	scp := scp.Server{}
	scp.DB = deps.DB
	scp.Logger = deps.Logger
	scp.CertPath = "certs/scp.rsa"
	scp.BindAddr = config.SCPBindAddr
	go scp.ListenAndServe()

	// Start server
	e.Run(s)
}