package handler

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/service"
	"github.com/pagient/pagient-server/internal/ui/renderer"
	"github.com/pagient/pagient-server/internal/ui/router/middleware/context"

	"github.com/go-chi/render"
)

// GetPatients lists all patients
func GetPatients(patientService service.PatientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		patients, err := patientService.ListPatients()
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewPatientListResponse(patients))
	}
}

// AddPatient adds a patient
func AddPatient(patientService service.PatientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		patientReq := &renderer.PatientRequest{}
		if err := render.Bind(req, patientReq); err != nil {
			render.Render(w, req, renderer.ErrBadRequest(err))
			return
		}

		// Set clientID to the client that added the patientReq
		ctxClient := req.Context().Value(context.ClientKey).(*model.Client)
		if ctxClient == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			return
		}
		patientReq.ClientID = ctxClient.ID

		patient := patientReq.GetModel()
		err := patientService.CreatePatient(patient)
		if err != nil {
			if service.IsModelExistErr(err) {
				render.Render(w, req, renderer.ErrConflict(err))
				return
			}

			if service.IsModelValidationErr(err) {
				render.Render(w, req, renderer.ErrValidation(err))
				return
			}

			// on any other error raise 500 status
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.Status(req, http.StatusCreated)
		render.Render(w, req, renderer.NewPatientResponse(patient))
	}
}

// GetPatient returns the patient by specified id
func GetPatient() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctxPatient := req.Context().Value(context.PatientKey).(*model.Patient)

		render.Render(w, req, renderer.NewPatientResponse(ctxPatient))
	}
}

// UpdatePatient updates a patient by specified id
func UpdatePatient(patientService service.PatientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		patientReq := &renderer.PatientRequest{}
		if err := render.Bind(req, patientReq); err != nil {
			render.Render(w, req, renderer.ErrBadRequest(err))
			return
		}

		// prevent ID update
		// prevent direct ClientID update
		ctxPatient := req.Context().Value(context.PatientKey).(*model.Patient)
		patientReq.ID = ctxPatient.ID
		patientReq.ClientID = ctxPatient.ClientID

		// Set clientID to the client that updated the patient
		// Update/Keep ClientID of requester's client
		ctxClient := req.Context().Value(context.ClientKey).(*model.Client)
		if ctxClient != nil {
			patientReq.ClientID = ctxClient.ID
		}

		patient := patientReq.GetModel()
		err := patientService.UpdatePatient(patient)
		if err != nil {
			if service.IsModelValidationErr(err) {
				render.Render(w, req, renderer.ErrValidation(err))
				return
			}

			if service.IsExternalServiceErr(err) {
				render.Render(w, req, renderer.ErrGateway(err))
				return
			}

			// on any other error raise 500 status
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.Render(w, req, renderer.NewPatientResponse(patient))
	}
}

// DeletePatient deletes a patient by specified id
func DeletePatient(patientService service.PatientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctxPatient := req.Context().Value(context.PatientKey).(*model.Patient)

		if err := patientService.DeletePatient(ctxPatient); err != nil {
			if service.IsInvalidArgumentErr(err) {
				render.Render(w, req, renderer.ErrBadRequest(err))
				return
			}

			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
