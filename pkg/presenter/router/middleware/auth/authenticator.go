package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/pagient/pagient-server/pkg/service"
)

// Authenticator middleware is used to authenticate the user by bearer token
func Authenticator(tokenService service.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token, claims, err := jwtauth.FromContext(req.Context())

			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if token == nil || !token.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			username, ok := claims.Get("user")
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokens, err := tokenService.Get(username.(string))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			for _, tok := range tokens {
				if tok.Token == token.Raw {
					// Token is authenticated, pass it through
					next.ServeHTTP(w, req)
					return
				}
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		})
	}
}
