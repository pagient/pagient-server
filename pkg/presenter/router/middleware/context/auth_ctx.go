package context

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/rs/zerolog/log"
)

// AuthCtx middleware is used to load a User object from
// the basic auth headers passed through as the request. In case
// the User could not be found, we stop here and return a 500.
func AuthCtx(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			_, claims, err := jwtauth.FromContext(req.Context())
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			username, ok := claims.Get("user")
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			user, err := userService.Get(username.(string))
			if err != nil {
				log.Error().
					Err(err).
					Msg("get user failed")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(req.Context(), "user", user)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
