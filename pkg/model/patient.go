package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
)

// PatientState hold the state of the Patient
type PatientState string

// enumerates all states a patient can be in
const (
	// PatientStateNew is for when the Patient is Pending
	PatientStatePending PatientState = "pending"
	// PatientStateCalled is for when the Patient's Pager has been called
	PatientStateCalled PatientState = "called"
)

// Patient struct
type Patient struct {
	ID       uuid.UUID    `json:"id"`
	Name     string       `json:"name"`
	PagerID  int          `json:"pager_id,omitempty"`
	ClientID int          `json:"client_id,omitempty"`
	Status   PatientState `json:"status"`
}

// Call calls the pager he is associated with
func (patient *Patient) Call() error {
	log.Debug().
		Str("patient", patient.Name).
		Msg("patient has been called")

	pager, err := GetPagerByID(patient.PagerID)
	if err != nil {
		return err
	}

	err = pager.Call()

	return err
}

// Validate validates the patient
func (patient *Patient) Validate() error {
	var pagerIDs []int
	if patient.PagerID != 0 {
		pagers, err := GetPagers()
		if err != nil {
			return err
		}
		for _, pager := range pagers {
			pagerIDs = append(pagerIDs, pager.ID)
		}
	}

	return validation.ValidateStruct(&patient,
		validation.Field(&patient.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&patient.PagerID, validation.In(pagerIDs)),
		validation.Field(&patient.Status, validation.In(PatientStatePending, PatientStateCalled)),
	)
}

// GetPatients lists all patients
func GetPatients() ([]*Patient, error) {
	return db.GetPatients()
}

// GetPatient returns a patient by ID
func GetPatient(patientID uuid.UUID) (*Patient, error) {
	return db.GetPatient(patientID)
}

// SavePatient stores the values in the database
func SavePatient(patient *Patient) error {
	return db.AddPatient(patient)
}

// UpdatePatient updates the values in the database
func UpdatePatient(patient *Patient) error {
	return db.UpdatePatient(patient)
}

// RemovePatient deletes the values from the database
func RemovePatient(patient *Patient) error {
	return db.RemovePatient(patient)
}
