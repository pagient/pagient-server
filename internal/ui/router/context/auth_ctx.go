package context

import (
	"context"
	"net/http"

	"github.com/pagient/pagient-server/internal/service"
	"github.com/pagient/pagient-server/internal/ui/renderer"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// AuthCtx middleware is used to load a User object from
// the authentication headers passed through as the request. In case
// the User could not be found, we stop here and return a 500.
func AuthCtx(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			jwtToken, _, err := jwtauth.FromContext(req.Context())
			if err != nil {
				render.Render(w, req, renderer.ErrUnauthorized)
				return
			}

			user, err := userService.ShowUserByToken(jwtToken.Raw)
			if err != nil {
				log.Error().
					Err(err).
					Msg("get user failed")

				render.Render(w, req, renderer.ErrInternalServer(err))
				return
			}

			ctx := context.WithValue(req.Context(), UserKey, user)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
