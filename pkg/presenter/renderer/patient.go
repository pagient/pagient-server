package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/presenter/router/middleware/context"
	"github.com/pkg/errors"
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
	var patient *model.Patient

	// Request is an update
	if r.Context().Value(context.PatientKey) != nil {
		patient = r.Context().Value(context.PatientKey).(*model.Patient)

		if pr.ID != 0 && pr.ID != patient.ID {
			return errors.New("id attribute is not allowed to be updated")
		}

		if pr.ClientID != 0 && pr.ClientID != patient.ClientID {
			return errors.New("client_id attribute is not allowed to be updated")
		}

		if pr.PagerID == 0 && pr.Status == string(model.PatientStateCall) {
			return errors.New("patient call state can only be set if a pager is assigned")
		}
	} else {
		if pr.ClientID != 0 {
			return errors.New("client_id not allowed")
		}

		if pr.Status != "" {
			return errors.New("status not allowed")
		}
	}

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
		Status:           model.PatientState(pr.Status),
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
