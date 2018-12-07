package service

import (
	"github.com/pagient/pagient-easy-call-go/easycall"
	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetAll returns all patients
func (service *DefaultService) ListPatients() ([]*model.Patient, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	patients, err := tx.GetPatients()
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "get all patients failed")
	}

	tx.Commit()
	return patients, nil
}

func (service *DefaultService) ListPagerPatientsByStatus(states ...model.PatientState) ([]*model.Patient, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	patients, err := tx.GetPatientsWithPagerByStatus(states...)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "get patients having pagers by status failed")
	}

	tx.Commit()
	return patients, nil
}

// Get returns a patient by it's id
func (service *DefaultService) ShowPatient(id uint) (*model.Patient, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	patient, err := tx.GetPatient(id)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "get patient failed")
	}

	tx.Commit()
	return patient, nil
}

// Add adds a new patient if given model is valid and not already existing
func (service *DefaultService) CreatePatient(patient *model.Patient) (*model.Patient, error) {
	patient.Status = model.PatientStatePending

	if patient.ClientID == 0 {
		return nil, &invalidArgumentErr{"clientId: cannot be blank"}
	}

	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	if err := service.validatePatient(tx, patient); err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	if patient.Active {
		if err := service.markPatientsInactiveFromClient(tx, patient.ClientID); err != nil {
			tx.Rollback()
			return nil, errors.WithStack(err)
		}
	}

	if err := service.removeInactivePatientsWithoutPagerFromClient(tx, patient.ClientID); err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	patient, err = tx.AddPatient(patient)
	if err != nil {
		tx.Rollback()

		if isEntryNotValidErr(err) {
			return nil, &modelValidationErr{err.Error()}
		}

		if isEntryExistErr(err) {
			return nil, &modelExistErr{"patient already exists"}
		}

		return nil, errors.Wrap(err, "add patient failed")
	}

	tx.Commit()
	service.notifier.NotifyNewPatient(patient)

	return patient, errors.Wrap(err, "add patient failed")
}

// Update updates an existing patient if given model is valid
func (service *DefaultService) UpdatePatient(patient *model.Patient) (*model.Patient, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	if err := service.validatePatient(tx, patient); err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	// load patient's old state to compare changed properties
	patientBeforeUpdate, err := tx.GetPatient(patient.ID)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "get patient failed")
	}

	if patient.Active {
		if err := service.markPatientsInactiveFromClient(tx, patient.ClientID); err != nil {
			tx.Rollback()
			return nil, errors.WithStack(err)
		}
	}

	patient, err = tx.UpdatePatient(patient)
	if err != nil {
		tx.Rollback()

		if isEntryNotValidErr(err) {
			return nil, &modelValidationErr{err.Error()}
		}

		if isEntryNotExistErr(err) {
			return nil, &modelNotExistErr{"patient doesn't exist"}
		}

		return nil, errors.Wrap(err, "update patient failed")
	}

	if err := service.removeInactivePatientsWithoutPagerFromClient(tx, patient.ClientID); err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	// RoomAssignment status changed from another state to PatientStateCall
	if patient.Status == model.PatientStateCall && patient.Status != patientBeforeUpdate.Status {
		log.Debug().
			Uint("pager", patient.PagerID).
			Msg("pager gets called")

		if err := service.callPatient(tx, patient); err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "call patient failed")
		}
	}

	tx.Commit()
	service.notifier.NotifyUpdatedPatient(patient)

	return patient, nil
}

// Remove deletes an existing patient
func (service *DefaultService) DeletePatient(patient *model.Patient) error {
	if patient.PagerID != 0 {
		return &invalidArgumentErr{"pagerId: cannot be set"}
	}

	tx, err := service.db.Begin()
	if err != nil {
		return errors.Wrap(err, "create transaction failed")
	}

	err = tx.RemovePatient(patient)
	if err != nil {
		tx.Rollback()

		if isEntryNotExistErr(err) {
			return &modelNotExistErr{"patient doesn't exist"}
		}

		return errors.Wrap(err, "remove patient failed")
	}

	tx.Commit()
	service.notifier.NotifyDeletedPatient(patient)

	return nil
}

func (service *DefaultService) CallPatient(patient *model.Patient) error {
	tx, err := service.db.Begin()
	if err != nil {
		return errors.Wrap(err, "create transaction failed")
	}

	if err := service.callPatient(tx, patient); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "call patient failed")
	}

	tx.Commit()
	service.notifier.NotifyUpdatedPatient(patient)

	return nil
}

func (service *DefaultService) callPatient(tx Tx, patient *model.Patient) error {
	client := easycall.NewClient(config.EasyCall.URL, config.EasyCall.User, config.EasyCall.Password)

	if err := client.Send(&easycall.SendOptions{
		Receiver: int(patient.PagerID),
		Message:  "",
		Port:     config.EasyCall.Port,
	}); err != nil {
		return &externalServiceErr{"pager call failed"}
	}

	patient.Status = model.PatientStateCalled

	patient, err := tx.UpdatePatient(patient)
	if err != nil {
		return errors.Wrap(err, "update patient failed")
	}

	return nil
}

func (service *DefaultService) validatePatient(tx Tx, patient *model.Patient) error {
	var pagers []*model.Pager

	if patient.PagerID != 0 {
		// load pagers to validate if pager sent with request is valid
		var err error
		pagers, err = tx.GetUnassignedPagers()
		if err != nil {
			return errors.Wrap(err, "get all pagers failed")
		}
	} else if patient.Status == model.PatientStateCall {
		return &modelValidationErr{"Status: \"call\" can only be set if PagerID is set."}
	}

	pat, err := tx.GetPatient(patient.ID)
	if err != nil {
		return errors.Wrap(err, "get patient failed")
	}

	// patient exists so it is an update
	// pager hasn't changed so it is also valid
	if pat != nil && pat.PagerID == patient.PagerID {
		pagers = append(pagers, &model.Pager{ID: patient.PagerID})
	}

	// validate patient
	if err := patient.Validate(pagers); err != nil {
		if model.IsValidationErr(err) {
			return &modelValidationErr{err.Error()}
		}

		return errors.Wrap(err, "validate patient failed")
	}

	return nil
}

func (service *DefaultService) markPatientsInactiveFromClient(tx PatientTx, clientID uint) error {
	patients, err := tx.GetPatientsByClient(clientID, true)
	if err != nil {
		return errors.Wrap(err, "get all patients by client failed")
	}

	if err := tx.MarkPatientsInactiveByClient(clientID); err != nil {
		return errors.Wrap(err, "mark all patients as inactive failed")
	}

	for _, patient := range patients {
		service.notifier.NotifyUpdatedPatient(patient)
	}

	return nil
}

func (service *DefaultService) removeInactivePatientsWithoutPagerFromClient(tx PatientTx, clientID uint) error {
	patients, err := tx.GetPatientsByClient(clientID, false, false)
	if err != nil {
		return errors.Wrap(err, "get all patients by client failed")
	}

	if err := tx.RemovePatientsByClient(clientID, false, false); err != nil {
		return errors.Wrap(err, "remove all inactive patients without pager failed")
	}

	for _, patient := range patients {
		service.notifier.NotifyDeletedPatient(patient)
	}

	return nil
}
