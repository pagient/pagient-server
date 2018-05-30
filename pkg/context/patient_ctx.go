package context

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/renderer"
	"github.com/satori/go.uuid"
)

// PatientCtx middleware is used to load a Patient object from
// the URL parameters passed through as the request. In case
// the Patient could not be found, we stop here and return a 404.
func PatientCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var patient *model.Patient
		var err error

		if patientID := chi.URLParam(req, "patientID"); patientID != "" {
			var id uuid.UUID
			id, err = uuid.FromString(patientID)

			if err == nil {
				patient, err = model.GetPatient(id)
			}
		} else {
			render.Render(w, req, renderer.ErrNotFound)
			return
		}

		if err != nil {
			render.Render(w, req, renderer.ErrNotFound)
			return
		}

		ctx := context.WithValue(req.Context(), "patient", patient)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
