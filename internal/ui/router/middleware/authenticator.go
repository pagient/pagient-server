package middleware

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/service"
	"github.com/pagient/pagient-server/internal/ui/renderer"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// Authenticator middleware is used to authenticate the user by bearer token
func Authenticator(tokenService service.TokenService, userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			jwtToken, _, err := jwtauth.FromContext(req.Context())

			if err != nil {
				render.Render(w, req, renderer.ErrUnauthorized)
				return
			}

			if jwtToken == nil || !jwtToken.Valid {
				render.Render(w, req, renderer.ErrUnauthorized)
				return
			}

			token, err := tokenService.ShowToken(jwtToken.Raw)
			if err != nil {
				render.Render(w, req, renderer.ErrInternalServer(err))
				return
			}

			if token != nil {
				// Token is authenticated, pass it through
				next.ServeHTTP(w, req)
				return
			}

			render.Render(w, req, renderer.ErrUnauthorized)
			return
		})
	}
}
