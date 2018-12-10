package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

// PatientStatus hold the state of the Patient
type PatientStatus string

// enumerates all states a patient can be in
const (
	// PatientStatusPending is for when the patient is pending
	PatientStatusPending PatientStatus = "pending"
	// PatientStatusCall is for when the patient's pager gets called
	PatientStatusCall PatientStatus = "call"
	// PatientStatusCalled is for when the patient's pager has been called
	PatientStatusCalled PatientStatus = "called"
	// PatientStatusFinished is for when the patient is finished with his medical examination
	PatientStatusFinished PatientStatus = "finished"
)

// Patient struct
type Patient struct {
	ID               uint   `gorm:"primary_key"`
	SocialSecurityNo string `gorm:"column:ssn;not null;unique"`
	Name             string `gorm:"not null"`
	Pager            Pager  `gorm:"save_associations:false"`
	PagerID          uint
	Client           Client `gorm:"save_associations:false"`
	ClientID         uint
	Status           PatientStatus `gorm:"not null" sql:"default:\"pending\""`
	Active           bool          `gorm:"not null" sql:"default:false"`
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
		validation.Field(&patient.SocialSecurityNo, validation.Required, is.Digit, validation.Length(10, 10)),
		validation.Field(&patient.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&patient.PagerID, validation.In(pagerIDs...)),
		validation.Field(&patient.Status, validation.In(PatientStatusPending, PatientStatusCall, PatientStatusCalled, PatientStatusFinished)),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occured")
		}

		return &modelValidationErr{err.Error()}
	}

	return nil
}
