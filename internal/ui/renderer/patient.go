package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/internal/model"
)

// PatientRequest is the request payload for patient data model
type PatientRequest struct {
	ID               uint   `json:"id"`
	SocialSecurityNo string `json:"ssn"`
	Name             string `json:"name"`
	PagerID          uint   `json:"pagerId"`
	ClientID         uint   `json:"clientId"`
	Status           string `json:"status"`
	Active           bool   `json:"active"`
}

// Bind postprocesses the decoding of the request body
func (pr *PatientRequest) Bind(r *http.Request) error {
	return nil
}

// GetModel returns a Patient model
func (pr *PatientRequest) GetModel() *model.Patient {
	return &model.Patient{
		ID:               pr.ID,
		SocialSecurityNo: pr.SocialSecurityNo,
		Name:             pr.Name,
		PagerID:          pr.PagerID,
		ClientID:         pr.ClientID,
		Status:           model.PatientStatus(pr.Status),
		Active:           pr.Active,
	}
}

// PatientResponse is the response payload for the patient data model
type PatientResponse struct {
	ID               uint   `json:"id"`
	SocialSecurityNo string `json:"ssn"`
	Name             string `json:"name"`
	PagerID          uint   `json:"pagerId,omitempty"`
	ClientID         uint   `json:"clientId"`
	Status           string `json:"status"`
	Active           bool   `json:"active"`
}

// NewPatientResponse creates a new patient response from patient model
func NewPatientResponse(patient *model.Patient) *PatientResponse {
	resp := &PatientResponse{
		ID:               patient.ID,
		SocialSecurityNo: patient.SocialSecurityNo,
		Name:             patient.Name,
		PagerID:          patient.PagerID,
		ClientID:         patient.ClientID,
		Status:           string(patient.Status),
		Active:           patient.Active,
	}

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
