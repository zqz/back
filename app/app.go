package app

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/zqzca/back/dependencies"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/ws"
	"golang.org/x/crypto/acme/autocert"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
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

	// Logging
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{}

	// Websockets
	ws := ws.NewServer()

	// Shared dependencies between all controller
	deps := dependencies.Dependencies{
		Fs:     afero.NewOsFs(),
		Logger: log,
		DB:     db,
		WS:     ws,
	}

	ws.Dependencies = &deps
	go ws.Start()

	// // Start SCP
	// scp := scp.Server{}
	// scp.DB = deps.DB
	// scp.Logger = deps.Logger
	// scp.CertPath = "certs/scp.rsa"
	// scp.BindAddr = config.SCPBindAddr
	// go scp.ListenAndServe()

	if config.Secure {
		c := autocert.DirCache("certs")
		m := autocert.Manager{
			Cache:      c,
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("x.zqz.ca", "de.zqz.ca", "zqz.ca"),
		}

		s := &http.Server{
			Addr:      config.HTTPBindAddr,
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
			Handler:   Routes(deps),
		}

		deps.Info("Listening for HTTP1.1 Connections", "addr", ":3001")
		deps.Info("Listening for HTTP2 Connections", "addr", config.HTTPBindAddr)
		go http.ListenAndServe(":3001", secureRedirect())
		s.ListenAndServeTLS("", "")
	} else {
		deps.Info("Listening for HTTP1.1 Connections", "addr", config.HTTPBindAddr)
		http.ListenAndServe(config.HTTPBindAddr, Routes(deps))
	}
}
