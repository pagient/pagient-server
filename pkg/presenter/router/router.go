package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/assets"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/presenter/handler"
	"github.com/pagient/pagient-api/pkg/presenter/router/middleware/basicauth"
	"github.com/pagient/pagient-api/pkg/presenter/router/middleware/context"
	"github.com/pagient/pagient-api/pkg/presenter/router/middleware/header"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// Load initializes the routing of the application.
func Load(cfg *config.Config, clientHandler *handler.ClientHandler, pagerHandler *handler.PagerHandler, patientHandler *handler.PatientHandler,
	websocketHandler *handler.WebsocketHandler, clientService service.ClientService, patientService service.PatientService) http.Handler {

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
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.Route("/", func(root chi.Router) {
		root.Use(basicauth.Basicauth(cfg))

		root.Route("/api", func(r chi.Router) {
			r.Use(render.SetContentType(render.ContentTypeJSON))
			r.Use(context.AuthCtx(clientService))

			// Manage patients
			r.Route("/patients", func(r chi.Router) {
				r.Get("/", patientHandler.GetPatients)
				r.Post("/", patientHandler.AddPatient)

				r.Route("/{patientID}", func(r chi.Router) {
					r.Use(context.PatientCtx(patientService))

					r.Get("/", patientHandler.GetPatient)
					r.Post("/", patientHandler.UpdatePatient)
					r.Delete("/", patientHandler.DeletePatient)
				})
			})

			// List pagers
			r.Get("/pagers", pagerHandler.GetPagers)
			// List clients
			r.Get("/clients", clientHandler.GetClients)
		})

		// Serve Websocket
		root.Get("/ws", websocketHandler.ServeWebsocket)

		// Pagient UI static files
		root.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// static files contain all files from "public/dist/"
			fs := http.StripPrefix("/", http.FileServer(assets.HTTP))

			fs.ServeHTTP(w, req)
		}))
	})

	return mux
}
