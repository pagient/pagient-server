package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/context"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/renderer"
	"github.com/pagient/pagient-api/pkg/websocket"
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
func AddPatient(cfg *config.Config, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := &renderer.PatientRequest{}
		if err := render.Bind(req, data); err != nil {
			render.Render(w, req, renderer.ErrBadRequest(err))
			return
		}

		patient := data.Patient

		// Set clientID to the client that added the patient
		ctxClient := req.Context().Value(context.ClientKey).(*model.Client)
		patient.ClientID = ctxClient.ID

		patient.Status = model.PatientStatePending

		if err := patient.Validate(); err != nil {
			render.Render(w, req, renderer.ErrValidation(err))
			return
		}

		if patient, _ := model.GetPatient(patient.ID); patient != nil {
			render.Render(w, req, renderer.ErrConflict(fmt.Errorf("patient with id %d already exists", patient.ID)))
			return
		}

		if err := model.SavePatient(patient); err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		// Broadcast new active patient in websocket hub
		if patient.Active {
			err := hub.Broadcast(websocket.MessageTypePatientUpdate, patient)
			if err != nil {
				render.Render(w, req, renderer.ErrInternalServer(err))
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
		ctxPatient := req.Context().Value(context.PatientKey).(*model.Patient)

		if err := render.Render(w, req, renderer.NewPatientResponse(ctxPatient)); err != nil {
			render.Render(w, req, renderer.ErrRender(err))
		}
	}
}

// UpdatePatient updates a patient by specified id
func UpdatePatient(cfg *config.Config, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := &renderer.PatientRequest{}
		if err := render.Bind(req, data); err != nil {
			render.Render(w, req, renderer.ErrBadRequest(err))
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

		// Broadcast patient status in websocket hub
		if err := hub.Broadcast(websocket.MessageTypePatientUpdate, patient); err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		if patient.Status == model.PatientStateCall {
			if err := patient.Call(); err != nil {
				render.Render(w, req, renderer.ErrGateway(err))
				return
			}
			patient.Status = model.PatientStateCalled

			// Broadcast patient status in websocket hub
			err := hub.Broadcast(websocket.MessageTypePatientUpdate, patient)
			if err != nil {
				render.Render(w, req, renderer.ErrInternalServer(err))
				return
			}
		}

		render.Render(w, req, renderer.NewPatientResponse(patient))
	}
}

// DeletePatient deletes a patient by specified id
func DeletePatient(cfg *config.Config, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctxPatient := req.Context().Value(context.PatientKey).(*model.Patient)

		if err := model.RemovePatient(ctxPatient); err != nil {
			render.Render(w, req, renderer.ErrBadRequest(err))
			return
		}

		// Broadcast remove patient in websocket hub
		err := hub.Broadcast(websocket.MessageTypePatientDelete, ctxPatient)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
