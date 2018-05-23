package handler

import (
	"net/http"

	"github.com/pagient/pagient-api/pkg/config"
)

// Notfound just returns a 404 not found error.
func Notfound(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}
