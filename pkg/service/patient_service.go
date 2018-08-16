package service

import (
	"github.com/pagient/pagient-easy-call-go/easycall"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/notifier"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// PatientService interface
type PatientService interface {
	GetAll() ([]*model.Patient, error)
	Get(uint) (*model.Patient, error)
	Add(*model.Patient) (*model.Patient, error)
	Update(*model.Patient) (*model.Patient, error)
	Remove(*model.Patient) error
}

// DefaultPatientService struct
type DefaultPatientService struct {
	cfg               *config.Config
	patientRepository PatientRepository
	pagerRepository   PagerRepository
	notifier          notifier.Notifier
}

// NewPatientService initializes a PatientService
func NewPatientService(cfg *config.Config, patientRepository PatientRepository, pagerRepository PagerRepository, notifier notifier.Notifier) PatientService {
	return &DefaultPatientService{
		cfg:               cfg,
		patientRepository: patientRepository,
		pagerRepository:   pagerRepository,
		notifier:          notifier,
	}
}

// GetAll returns all patients
func (service *DefaultPatientService) GetAll() ([]*model.Patient, error) {
	session := service.patientRepository.BeginTx()
	patients, err := service.patientRepository.GetAll(session)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all patients failed")

		service.patientRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get all patients failed")
	}

	service.patientRepository.CommitTx(session)
	return patients, nil
}

// Get returns a patient by it's id
func (service *DefaultPatientService) Get(id uint) (*model.Patient, error) {
	session := service.patientRepository.BeginTx()
	patient, err := service.patientRepository.Get(session, id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("patient ID", id).
			Msg("get patient failed")

		service.patientRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get patient failed")
	}

	service.patientRepository.CommitTx(session)
	return patient, nil
}

// Add adds a new patient if given model is valid and not already existing
func (service *DefaultPatientService) Add(patient *model.Patient) (*model.Patient, error) {
	patient.Status = model.PatientStatePending

	if patient.ClientID == 0 {
		return nil, &invalidArgumentErr{"clientId: cannot be blank"}
	}

	session := service.patientRepository.BeginTx()
	if err := service.validatePatient(session, patient); err != nil {
		return nil, errors.WithStack(err)
	}

	patient, err := service.patientRepository.Add(session, patient)
	if err != nil {
		service.patientRepository.RollbackTx(session)

		if isEntryNotValidErr(err) {
			return nil, &modelValidationErr{err.Error()}
		}

		if isEntryExistErr(err) {
			return nil, &modelExistErr{"patient already exists"}
		}

		log.Error().
			Err(err).
			Msg("add patient failed")

		return nil, errors.Wrap(err, "add patient failed")
	}

	if patient.Active {
		if err := service.cleanupPatients(session, patient); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	service.patientRepository.CommitTx(session)
	service.notifier.NotifyNewPatient(patient)

	return patient, errors.Wrap(err, "add patient failed")
}

// Update updates an existing patient if given model is valid
func (service *DefaultPatientService) Update(patient *model.Patient) (*model.Patient, error) {
	session := service.patientRepository.BeginTx()
	if err := service.validatePatient(session, patient); err != nil {
		return nil, errors.WithStack(err)
	}

	// load patient's old state to compare changed properties
	patientBeforeUpdate, err := service.patientRepository.Get(session, patient.ID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("patient ID", patient.ID).
			Msg("get patient failed")

		service.patientRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get patient failed")
	}

	patient, err = service.patientRepository.Update(session, patient)
	if err != nil {
		service.patientRepository.RollbackTx(session)

		if isEntryNotValidErr(err) {
			return nil, &modelValidationErr{err.Error()}
		}

		if isEntryNotExistErr(err) {
			return nil, &modelNotExistErr{"patient doesn't exist"}
		}

		log.Error().
			Err(err).
			Msg("update patient failed")

		return nil, errors.Wrap(err, "update patient failed")
	}

	if patient.Active || (patient.PagerID == 0 && patient.PagerID != patientBeforeUpdate.PagerID) {
		if err := service.cleanupPatients(session, patient); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Patient status changed from another state to PatientStateCall
	if patient.Status == model.PatientStateCall && patient.Status != patientBeforeUpdate.Status {
		log.Debug().
			Uint("pager", patient.PagerID).
			Msg("pager gets called")

		client := easycall.NewClient(service.cfg.EasyCall.URL, service.cfg.EasyCall.User, service.cfg.EasyCall.Password)

		if err := client.Send(&easycall.SendOptions{
			Receiver: int(patient.PagerID),
			Message:  "",
			Port:     service.cfg.EasyCall.Port,
		}); err != nil {
			log.Error().
				Err(err).
				Uint("patient ID", patient.ID).
				Uint("pager ID", patient.PagerID).
				Msg("call pager failed")

			patient.Status = model.PatientStatePending

			patient, err = service.patientRepository.Update(session, patient)
			if err != nil {
				log.Error().
					Err(err).
					Msg("update patient failed")

				service.patientRepository.RollbackTx(session)
				return nil, errors.Wrap(err, "update patient failed")
			}

			return nil, &externalServiceErr{"pager call failed"}
		}

		patient.Status = model.PatientStateCalled

		patient, err = service.patientRepository.Update(session, patient)
		if err != nil {
			log.Error().
				Err(err).
				Msg("update patient failed")

			service.patientRepository.RollbackTx(session)
			return nil, errors.Wrap(err, "update patient failed")
		}
	}

	service.patientRepository.CommitTx(session)
	service.notifier.NotifyUpdatedPatient(patient)

	return patient, nil
}

// Remove deletes an existing patient
func (service *DefaultPatientService) Remove(patient *model.Patient) error {
	if patient.PagerID != 0 {
		return &invalidArgumentErr{"pagerId: cannot be set"}
	}

	session := service.patientRepository.BeginTx()
	patient, err := service.patientRepository.Remove(session, patient)
	if err != nil {
		service.patientRepository.RollbackTx(session)

		if isEntryNotExistErr(err) {
			return &modelNotExistErr{"patient doesn't exist"}
		}

		log.Error().
			Err(err).
			Msg("remove patient failed")

		return errors.Wrap(err, "remove patient failed")
	}

	service.patientRepository.CommitTx(session)
	service.notifier.NotifyDeletedPatient(patient)

	return nil
}

func (service *DefaultPatientService) validatePatient(session DB, patient *model.Patient) error {
	// load pagers to validate if pager sent with request is valid
	pagers, err := service.pagerRepository.GetUnassigned(session)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all pagers failed")

		return errors.Wrap(err, "get all pagers failed")
	}

	pat, err := service.patientRepository.Get(session, patient.ID)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get patient failed")
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

		log.Error().
			Err(err).
			Msg("validate patient failed")

		return errors.Wrap(err, "validate patient failed")
	}

	return nil
}

func (service *DefaultPatientService) cleanupPatients(session DB, patient *model.Patient) error {
	// mark all patients as inactive if current patient is active
	if patient.Active {
		updatedPatients, err := service.patientRepository.MarkAllExceptPatientInactiveByPatientClient(session, patient)
		if err != nil {
			log.Error().
				Err(err).
				Msg("mark all patients as inactive failed")

			service.patientRepository.RollbackTx(session)
			return errors.Wrap(err, "mark all patients as inactive failed")
		}

		for _, patient := range updatedPatients {
			service.notifier.NotifyUpdatedPatient(patient)
		}
	}

	// remove all inactive patients that have no pager assigned
	deletedPatients, err := service.patientRepository.RemoveAllExceptPatientInactiveNoPagerByPatientClient(session, patient)
	if err != nil {
		log.Error().
			Err(err).
			Msg("remove all inactive patients without pager failed")

		service.patientRepository.RollbackTx(session)
		return errors.Wrap(err, "remove all inactive patients without pager failed")
	}

	for _, patient := range deletedPatients {
		service.notifier.NotifyDeletedPatient(patient)
	}

	return nil
}
