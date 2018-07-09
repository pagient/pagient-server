package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pkg/errors"
)

// PatientRequest is the request payload for patient data model
type PatientRequest struct {
	*model.Patient
}

// Bind postprocesses the decoding of the request body
func (pr *PatientRequest) Bind(r *http.Request) error {
	var patient *model.Patient

	// Request is an update
	if r.Context().Value("patient") != nil {
		patient = r.Context().Value("patient").(*model.Patient)

		if pr.Patient.ID != patient.ID {
			return errors.New("id attribute is not allowed to be updated")
		}

		if pr.Patient.PagerID == 0 && pr.Patient.Status == model.PatientStateCall {
			return errors.New("patient call state can only be set if a pager is assigned")
		}
	} else {
		if pr.Patient.Status != "" {
			return errors.New("status not allowed")
		}
	}

	if pr.Patient.ClientID != 0 {
		return errors.New("client_id not allowed")
	}

	return nil
}

// PatientResponse is the response payload for the patient data model
type PatientResponse struct {
	*model.Patient
}

// NewPatientResponse creates a new patient response from patient model
func NewPatientResponse(patient *model.Patient) *PatientResponse {
	resp := &PatientResponse{Patient: patient}

	return resp
}

// Render preprocesses the response before marshalling
func (pr *PatientResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PatientListResponse is the list response payload for the patient data model
type PatientListResponse []*PatientResponse

// NewPatientListResponse creates a new patient list response from multiple patient models
func NewPatientListResponse(patients []*model.Patient) []render.Renderer {
	list := make([]render.Renderer, len(patients))
	for i, patient := range patients {
		list[i] = NewPatientResponse(patient)
	}
	return list
}
