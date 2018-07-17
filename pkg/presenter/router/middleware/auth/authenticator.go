package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/pagient/pagient-server/pkg/service"
)

// Authenticator middleware is used to authenticate the user by bearer token
func Authenticator(tokenService service.TokenService, userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			jwtToken, _, err := jwtauth.FromContext(req.Context())

			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if jwtToken == nil || !jwtToken.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			token, err := tokenService.Get(jwtToken.Raw)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if token != nil {
				// Token is authenticated, pass it through
				next.ServeHTTP(w, req)
				return
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		})
	}
}
