package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/renderer"
)

// GetPatients lists all patients
func GetPatients(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		patients, err := model.GetPatients()
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewPatientListResponse(patients))
	}
}

// AddPatient adds a patient
func AddPatient(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := &renderer.PatientRequest{}
		if err := render.Bind(req, data); err != nil {
			render.Render(w, req, renderer.ErrInvalidRequest(err))
			return
		}

		patient := data.Patient
		if patient.Status == "" {
			patient.Status = model.PatientStatePending
		}

		if err := patient.Validate(); err != nil {
			render.Render(w, req, renderer.ErrValidation(err))
			return
		}

		if err := model.SavePatient(patient); err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		if patient.Status == model.PatientStateCalled {
			if err := patient.Call(); err != nil {
				render.Render(w, req, renderer.ErrGateway(err))
				return
			}
		}

		render.Status(req, http.StatusCreated)
		render.Render(w, req, renderer.NewPatientResponse(patient))
	}
}

// GetPatient returns tha patient by specified id
func GetPatient(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctxPatient := req.Context().Value("patient").(*model.Patient)

		if err := render.Render(w, req, renderer.NewPatientResponse(ctxPatient)); err != nil {
			render.Render(w, req, renderer.ErrRender(err))
		}
	}
}

// UpdatePatient updates a patient by specified id
func UpdatePatient(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctxPatient := *req.Context().Value("patient").(*model.Patient)

		data := &renderer.PatientRequest{Patient: &ctxPatient}
		if err := render.Bind(req, data); err != nil {
			render.Render(w, req, renderer.ErrInvalidRequest(err))
			return
		}

		patient := data.Patient

		if err := patient.Validate(); err != nil {
			render.Render(w, req, renderer.ErrValidation(err))
			return
		}

		if err := model.UpdatePatient(patient); err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		if patient.Status == model.PatientStateCalled {
			if err := patient.Call(); err != nil {
				render.Render(w, req, renderer.ErrGateway(err))
				return
			}
		}

		render.Render(w, req, renderer.NewPatientResponse(patient))
	}
}

// DeletePatient deletes a patient by specified id
func DeletePatient(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctxPatient := req.Context().Value("patient").(*model.Patient)

		if err := model.RemovePatient(ctxPatient); err != nil {
			render.Render(w, req, renderer.ErrInvalidRequest(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
