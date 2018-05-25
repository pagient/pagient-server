package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/pagient/pagient-api/pkg/config"
)

// Basicauth integrates a simple basic authentication.
func Basicauth(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(cfg.General.Users) > 0 {
				w.Header().Set("WWW-Authenticate", `Basic realm="Pagient"`)

				s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

				if len(s) != 2 {
					http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
					return
				}

				b, err := base64.StdEncoding.DecodeString(s[1])

				if err != nil {
					http.Error(w, err.Error(), 401)
					return
				}

				pair := strings.SplitN(string(b), ":", 2)

				if len(pair) != 2 {
					http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
					return
				}

				pw, err := cfg.General.GetPassword(pair[0])
				if err != nil || pair[1] != pw {
					http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
