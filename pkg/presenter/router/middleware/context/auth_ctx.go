package context

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/presenter/renderer"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/rs/zerolog/log"
)

// AuthCtx middleware is used to load a Client object from
// the basic auth headers passed through as the request. In case
// the Client could not be found, we stop here and return a 500.
func AuthCtx(clientService service.ClientService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var client *model.Client
			var err error

			username, _, ok := r.BasicAuth()
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			}

			client, err = clientService.GetByName(username)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("get client failed")

				render.Render(w, r, renderer.ErrInternalServer(err))
				return
			}

			ctx := context.WithValue(r.Context(), "client", client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
