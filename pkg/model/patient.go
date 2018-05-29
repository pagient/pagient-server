package model

import (
	"github.com/satori/go.uuid"
	"github.com/rs/zerolog/log"
)

// PatientState hold the state of the Patient
type PatientState string

// enumerates all states a patient can be in
const (
	// PatientStateNew is for when the Patient is Pending
	PatientStatePending PatientState = "pending"
	// PatientStateAway is for when the Patient is Away and has the Pager
	PatientStateAway PatientState = "away"
	// PatientStateCalled is for when the Patient's Pager has been called
	PatientStateCalled PatientState = "called"
)

// Patient struct
type Patient struct {
	ID       uuid.UUID
	Name     string
	PagerID  int
	ClientID int
	Status   PatientState
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

// GetPatients lists all patients
func GetPatients() ([]*Patient, error) {
	return db.GetPatients()
}

// GetPatient returns a patient by ID
func GetPatient(patientID uuid.UUID) (*Patient, error) {
	return db.GetPatient(patientID)
}

// StorePatient stores the values in the database
func StorePatient(patient *Patient) error {
	if patient.ID != uuid.Nil {
		if err := db.UpdatePatient(patient); err != nil {
			return err
		}
	}

	return db.AddPatient(patient)
}

// RemovePatient deletes the values from the database
func RemovePatient(patient *Patient) error {
	return db.RemovePatient(patient)
}
