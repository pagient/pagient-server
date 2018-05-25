package model

import "github.com/rs/zerolog/log"

// PatientState hold the state of the Patient
type PatientState string

const (
	// PatientStateNew is for when the Patient is Pending
	PatientStatePending PatientState = "pending"
	// PatientStateAway is for when the Patient is Away and has the Pager
	PatientStateAway    PatientState = "away"
	// PatientStateCalled is for when the Patient's Pager has been called
	PatientStateCalled  PatientState = "called"
)

type Patient struct {
	ID       int64
	Name     string
	Pager    *Pager
	ClientID int64
	Status   PatientState
}

func (patient *Patient) Call() {
	log.Debug().
		Str("patient", patient.Name).
		Msg("patient has been called")

	patient.Pager.Call()
}

func GetPatients() ([]*Patient, error) {
	return db.GetPatients()
}

func GetPatient(patientID int64) (*Patient, error) {
	return db.GetPatient(patientID)
}

func StorePatient(patient *Patient) error {
	return db.AddPatient(patient)
}

func RemovePatient(patient *Patient) error {
	return db.RemovePatient(patient)
}
