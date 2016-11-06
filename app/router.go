package app

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/zqzca/back/controller"
	"github.com/zqzca/back/controller/chunks"
	"github.com/zqzca/back/controller/dashboard"
	"github.com/zqzca/back/controller/files"
	"github.com/zqzca/back/controller/thumbnails"
	"github.com/zqzca/back/dependencies"
	"github.com/zqzca/back/ws"
)

// Routes defines all the routes for the application
func Routes(deps dependencies.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	ep := deps.WS.(*ws.Server).Endpoint()
	r.Handle("/ws", ep)

	// Default
	r.Group(func(r chi.Router) {
		r.Use(middleware.CloseNotify)
		r.Use(middleware.Compress(5))

		r.(*chi.Mux).FileServer("/assets", http.Dir("./assets"))

		files := files.Controller{Dependencies: deps}
		r.Get("/d/:slug", files.Download) // Short DL URL

		// Chunks
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/check/:hash", files.Status) // I dont like this URL

			chunks := chunks.NewController(deps)
			r.Route("/chunks", func(r chi.Router) {
				r.Post("/", chunks.Create)
			})

			r.Route("/files", func(r chi.Router) {
				r.Post("/", files.Create)
				r.With(controller.Pagination).Get("/", files.Index)
				r.Get("/:slug", files.Show)
				r.Get("/:slug/data", files.Download)
				r.Delete("/:slug/delete", files.Delete)
			})
			// r.Get("/files/:slug/process", files.Process)

			thumbnails := thumbnails.Controller{Dependencies: deps}
			r.Get("/thumbnails/:id", thumbnails.Download)

			dash := dashboard.Controller{Dependencies: deps}
			r.With(controller.Pagination).Get("/dashboard", dash.Index)
		})

		// Catchall
		r.Get("/*", Index)
	})

	// r.Use(middleware.Timeout(60 * time.Second))

	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	// }))

	// // Users
	// users := users.Controller{Dependencies: deps}
	// v1.Post("/users", users.Create)
	// v1.Get("/username/valid", users.ValidateUsername)
	// v1.Get("/users/:id", users.Read)

	// // Sessions
	// sessions := sessions.Controller{Dependencies: deps}
	// v1.Post("/sessions", sessions.Create)

	// // P2P
	// v1.Get("/p2p/signaling", standard.WrapHandler(http.HandlerFunc(p2p.Signaling())))
	// v1.Get("/p2p/:id", p2p.Join)
	// v1.Post("/p2p/:id", p2p.Answer)

	return r
}

func secureRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "https://"+req.Host+req.RequestURI, http.StatusMovedPermanently)
	}
}
