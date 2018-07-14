package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/presenter/websocket"
	"github.com/pagient/pagient-server/pkg/service"
)

// PatientHandler struct
type PatientHandler struct {
	patientService service.PatientService
	wsHub          *websocket.Hub
}

// NewPatientHandler initializes a PatientHandler
func NewPatientHandler(patientService service.PatientService, hub *websocket.Hub) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
		wsHub:          hub,
	}
}

// GetPatients lists all patients
func (handler *PatientHandler) GetPatients(w http.ResponseWriter, req *http.Request) {
	patients, err := handler.patientService.GetAll()
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	render.RenderList(w, req, renderer.NewPatientListResponse(patients))
}

// AddPatient adds a patient
func (handler *PatientHandler) AddPatient(w http.ResponseWriter, req *http.Request) {
	data := &renderer.PatientRequest{}
	if err := render.Bind(req, data); err != nil {
		render.Render(w, req, renderer.ErrBadRequest(err))
		return
	}

	patient := data.Patient

	// Set clientID to the client that added the patient
	ctxClient := req.Context().Value("client").(*model.Client)
	if ctxClient == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
		return
	}
	patient.ClientID = ctxClient.ID

	patient, err := handler.patientService.Add(patient)
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

// GetPatient returns the patient by specified id
func (handler *PatientHandler) GetPatient(w http.ResponseWriter, req *http.Request) {
	ctxPatient := req.Context().Value("patient").(*model.Patient)

	render.Render(w, req, renderer.NewPatientResponse(ctxPatient))
}

// UpdatePatient updates a patient by specified id
func (handler *PatientHandler) UpdatePatient(w http.ResponseWriter, req *http.Request) {
	data := &renderer.PatientRequest{}
	if err := render.Bind(req, data); err != nil {
		render.Render(w, req, renderer.ErrBadRequest(err))
		return
	}

	patient := data.Patient

	// Set clientID to the client that updated the patient
	ctxClient := req.Context().Value("client").(*model.Client)
	if ctxClient != nil {
		patient.ClientID = ctxClient.ID
	}

	patient, err := handler.patientService.Update(patient)
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

// DeletePatient deletes a patient by specified id
func (handler *PatientHandler) DeletePatient(w http.ResponseWriter, req *http.Request) {
	ctxPatient := req.Context().Value("patient").(*model.Patient)

	if err := handler.patientService.Remove(ctxPatient); err != nil {
		if service.IsInvalidArgumentErr(err) {
			render.Render(w, req, renderer.ErrBadRequest(err))
			return
		}

		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
