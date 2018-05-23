package handler

import (
	"net/http"

	"github.com/pagient/pagient-api/pkg/config"
)

func UpdatePatient(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
