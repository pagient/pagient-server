package handler

import (
	"github.com/pagient/pagient-api/pkg/config"
	"net/http"
)

func GetClients(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
