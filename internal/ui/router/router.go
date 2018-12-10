package router

import (
	"net/http"
	"time"

	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/service"
	"github.com/pagient/pagient-server/internal/ui/handler"
	"github.com/pagient/pagient-server/internal/ui/router/middleware/auth"
	"github.com/pagient/pagient-server/internal/ui/router/middleware/context"
	"github.com/pagient/pagient-server/internal/ui/router/middleware/header"
	"github.com/pagient/pagient-server/internal/ui/static"
	"github.com/pagient/pagient-server/internal/ui/websocket"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// Load initializes the routing of the application.
func Load(s service.Service, wsHub *websocket.Hub) http.Handler {

	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure())
	mux.Use(header.Options())

	tokenAuth := jwtauth.New("HS256", []byte(config.General.Secret), nil)

	mux.Route("/", func(root chi.Router) {
		root.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(auth.Authenticator(s, s))

			r.Route("/api", func(r chi.Router) {
				r.Use(render.SetContentType(render.ContentTypeJSON))
				r.Use(context.AuthCtx(s))

				// Manage patients
				r.Route("/patients", func(r chi.Router) {
					r.Get("/", handler.GetPatients(s))
					r.With(context.ClientCtx(s)).Post("/", handler.AddPatient(s))

					r.Route("/{patientID}", func(r chi.Router) {
						r.Use(context.PatientCtx(s))

						r.Get("/", handler.GetPatient())
						r.With(context.ClientCtx(s)).Post("/", handler.UpdatePatient(s))
						r.Delete("/", handler.DeletePatient(s))
					})
				})

				// List pagers
				r.Get("/pagers", handler.GetPagers(s))
				// List clients
				r.Get("/clients", handler.GetClients(s))
			})

			// Serve Websocket
			r.Get("/ws", handler.ServeWebsocket(s, wsHub))
		})

		root.Route("/oauth", func(r chi.Router) {
			r.Post("/token", handler.CreateToken(s, s))

			r.Route("/", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(auth.Authenticator(s, s))

				r.Delete("/token", handler.DeleteToken(s, wsHub))
				r.Get("/sessions", handler.GetSessions(s, s))
			})
		})

		// Pagient UI static files
		root.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// static files contain all files from "public/dist/"
			fs := http.StripPrefix("/", http.FileServer(static.HTTP))

			fs.ServeHTTP(w, req)
		}))
	})

	return mux
}
