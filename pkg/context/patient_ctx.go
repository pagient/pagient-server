package context

import (
	"net/http"
)

// PatientCtx loads specified patient in url into the context
func PatientCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
	})
}
