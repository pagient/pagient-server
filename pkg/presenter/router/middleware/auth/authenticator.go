package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/pagient/pagient-api/pkg/service"
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

			tokens, err := tokenService.Get(username.(string))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
				return
			}

			for _, tok := range tokens {
				if tok.Token == token.Raw {
					// Token is authenticated, pass it through
					next.ServeHTTP(w, req)
					return
				}
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			return
		})
	}
}
