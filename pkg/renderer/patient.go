package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/satori/go.uuid"
	"fmt"
)

// PatientRequest is the request payload for Patient data model.
type PatientRequest struct {
	*model.Patient
}

func (pr *PatientRequest) Bind(r *http.Request) error {
	var patient *model.Patient
	if r.Context().Value("patient") != nil {
		patient = r.Context().Value("patient").(*model.Patient)

		if pr.Patient.ID != patient.ID {
			return fmt.Errorf("id not allowed")
		}
	} else if pr.Patient.ID != uuid.Nil {
		return fmt.Errorf("id not allowed")
	}

	if pr.Patient.ClientID != 0 {
		return fmt.Errorf("client_id not allowed")
	}
	return nil
}

// PatientResponse is the response payload for the Patient data model.
type PatientResponse struct {
	*model.Patient
}

func NewPatientResponse(patient *model.Patient) *PatientResponse {
	resp := &PatientResponse{Patient: patient}

	return resp
}

func (pr *PatientResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

type PatientListResponse []*PatientResponse

func NewPatientListResponse(patients []*model.Patient) []render.Renderer {
	list := []render.Renderer{}
	for _, patient := range patients {
		list = append(list, NewPatientResponse(patient))
	}
	return list
}
