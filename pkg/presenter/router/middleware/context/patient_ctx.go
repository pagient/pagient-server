package context

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/presenter/renderer"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// PatientCtx middleware is used to load a Patient object from
// the URL parameters passed through as the request. In case
// the Patient could not be found, we stop here and return a 404.
func PatientCtx(patientService service.PatientService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var patient *model.Patient

			if patientID := chi.URLParam(req, "patientID"); patientID != "" {
				id, err := strconv.Atoi(patientID)
				if err != nil {
					render.Render(w, req, renderer.ErrBadRequest(err))
					return
				}

				patient, err = patientService.Get(id)
				if err != nil {
					log.Fatal().
						Err(err).
						Msg("get patient failed")

					render.Render(w, req, renderer.ErrInternalServer(err))
					return
				}

				if patient == nil {
					render.Render(w, req, renderer.ErrNotFound)
					return
				}

				ctx := context.WithValue(req.Context(), "patient", patient)
				next.ServeHTTP(w, req.WithContext(ctx))
				return
			}

			err := errors.New("patient id parameter missing in url")
			log.Fatal().
				Err(err).
				Msg("patient id parameter missing in url")

			render.Render(w, req, renderer.ErrInternalServer(err))
		})
	}
}
