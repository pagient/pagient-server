package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

// PatientState hold the state of the Patient
type PatientState string

// enumerates all states a patient can be in
const (
	// PatientStateNew is for when the Patient is Pending
	PatientStatePending PatientState = "pending"
	// PatientStateCall is for when the Patient's Pager gets called
	PatientStateCall PatientState = "call"
	// PatientStateCalled is for when the Patient's Pager has been called
	PatientStateCalled PatientState = "called"
	// PatientStateFinished is for when the Patient is Finished with his medical examination
	PatientStateFinished PatientState = "finished"
)

// Patient struct
type Patient struct {
	ID       int          `json:"id"`
	Ssn      string       `json:"ssn"`
	Name     string       `json:"name"`
	PagerID  int          `json:"pagerId,omitempty"`
	ClientID int          `json:"clientId,omitempty"`
	Status   PatientState `json:"status"`
	Active   bool         `json:"active"`
}

// Validate validates the patient
func (patient *Patient) Validate(pagers []*Pager) error {
	// convert pager slice to generic interface slice
	pagerIDs := make([]interface{}, len(pagers))
	for i, pager := range pagers {
		pagerIDs[i] = pager.ID
	}

	if err := validation.ValidateStruct(patient,
		validation.Field(&patient.ID, validation.Required),
		validation.Field(&patient.Ssn, validation.Required, is.Digit, validation.Length(10, 10)),
		validation.Field(&patient.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&patient.PagerID, validation.In(pagerIDs...)),
		validation.Field(&patient.Status, validation.In(PatientStatePending, PatientStateCall, PatientStateCalled, PatientStateFinished)),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occured")
		}

		return &modelValidationErr{err.Error()}
	}

	return nil
}
