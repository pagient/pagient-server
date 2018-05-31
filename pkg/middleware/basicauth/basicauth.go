package basicauth

import (
	"net/http"

	"github.com/pagient/pagient-api/pkg/config"
)

// Basicauth integrates a simple basic authentication.
func Basicauth(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(cfg.General.Users) > 0 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Pagient"`)

				username, password, ok := r.BasicAuth()
				if !ok {
					http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
					return
				}

				pw, err := cfg.General.GetPassword(username)
				if err != nil || password != pw {
					http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
