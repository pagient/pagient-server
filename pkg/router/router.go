package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/context"
	"github.com/pagient/pagient-api/pkg/handler"
	"github.com/pagient/pagient-api/pkg/middleware/basicauth"
	"github.com/pagient/pagient-api/pkg/middleware/header"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// Load initializes the routing of the application.
func Load(cfg *config.Config) http.Handler {
	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.NotFound(handler.Notfound(cfg))

	mux.Route("/", func(root chi.Router) {
		root.Use(basicauth.Basicauth(cfg))

		// Manage patients
		root.Route("/patients", func(state chi.Router) {
			state.Get("/", handler.GetPatients(cfg))
			state.Post("/", handler.AddPatient(cfg))

			state.Route("/{patientID}", func(state chi.Router) {
				state.Use(context.PatientCtx)

				state.Get("/", handler.GetPatient(cfg))
				state.Post("/", handler.UpdatePatient(cfg))
			})
		})

		// List pagers
		root.Get("/pagers", handler.GetPagers(cfg))
		// List clients
		root.Get("/clients", handler.GetClients(cfg))
	})

	return mux
}
