package context

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
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
					http.Error(w, "patient id not an integer", 400)
					return
				}

				patient, err = patientService.Get(id)
				if err != nil {
					log.Fatal().
						Err(err).
						Msg("get patient failed")

					http.Error(w, http.StatusText(500), 500)
					return
				}

				if patient == nil {
					http.Error(w, http.StatusText(404), 404)
					return
				}

				ctx := context.WithValue(req.Context(), "patient", patient)
				next.ServeHTTP(w, req.WithContext(ctx))
				return
			}

			err := errors.New("patient id parameter missing in url")
			log.Error().
				Err(err).
				Msg("patient id parameter missing in url")

			http.Error(w, http.StatusText(500), 500)
		})
	}
}
