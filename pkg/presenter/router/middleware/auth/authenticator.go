package auth

import (
	"net/http"

	"github.com/pagient/pagient-api/pkg/service"
	"github.com/go-chi/jwtauth"
)

func Authenticator(tokenService service.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token, claims, err := jwtauth.FromContext(req.Context())

			if err != nil {
				http.Error(w, http.StatusText(401), 401)
				return
			}

			if token == nil || !token.Valid {
				http.Error(w, http.StatusText(401), 401)
				return
			}

			username, ok := claims.Get("user")
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
				return
			}

			invalidToken, err := tokenService.Get(username.(string))
			if err != nil || invalidToken != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, req)
		})
	}
}
