package service

import (
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-easy-call-go/easycall"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// PatientService interface
type PatientService interface {
	GetAll() ([]*model.Patient, error)
	Get(int) (*model.Patient, error)
	Add(*model.Patient) error
	Update(*model.Patient) error
	Remove(*model.Patient) error
}

// DefaultPatientService struct
type DefaultPatientService struct {
	cfg               *config.Config
	patientRepository PatientRepository
	pagerRepository   PagerRepository
}

// NewPatientService initializes a PatientService
func NewPatientService(cfg *config.Config, patientRepository PatientRepository, pagerRepository PagerRepository) PatientService {
	return &DefaultPatientService{
		cfg:               cfg,
		patientRepository: patientRepository,
		pagerRepository:   pagerRepository,
	}
}

// GetAll returns all patients
func (service *DefaultPatientService) GetAll() ([]*model.Patient, error) {
	patients, err := service.patientRepository.GetAll()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all patients failed")
	}

	return patients, errors.Wrap(err, "get all patients failed")
}

// Get returns a patient by it's id
func (service *DefaultPatientService) Get(id int) (*model.Patient, error) {
	patient, err := service.patientRepository.Get(id)
	if err != nil {
		log.Error().
			Err(err).
			Int("patient ID", id).
			Msg("get patient failed")
	}

	return patient, errors.Wrap(err, "get patient failed")
}

// Add adds a new patient if given model is valid and not already existing
func (service *DefaultPatientService) Add(patient *model.Patient) error {
	patient.Status = model.PatientStatePending

	pagers, err := service.pagerRepository.GetAll()
	if err != nil {
		return errors.Wrap(err, "get all pagers failed")
	}

	if err := patient.Validate(pagers); err != nil {
		if model.IsValidationErr(err) {
			return errors.Wrap(err, "patient not valid")
		}

		log.Error().
			Err(err).
			Msg("validate patient failed")

		return errors.Wrap(err, "validate patient failed")
	}

	err = service.patientRepository.Add(patient)
	if err != nil {
		if isEntryNotValidErr(err) {
			return &modelValidationErr{err.Error()}
		}

		if isEntryExistErr(err) {
			return &modelExistErr{"patient already exists"}
		}

		log.Error().
			Err(err).
			Msg("add patient failed")
	}

	return errors.Wrap(err, "add patient failed")
}

// Update updates an existing patient if given model is valid
func (service *DefaultPatientService) Update(patient *model.Patient) error {
	pagers, err := service.pagerRepository.GetAll()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all pagers failed")

		return errors.Wrap(err, "get all pagers failed")
	}

	if err := patient.Validate(pagers); err != nil {
		if model.IsValidationErr(err) {
			return errors.Wrap(err, "patient not valid")
		}

		log.Error().
			Err(err).
			Msg("validate patient failed")

		return errors.Wrap(err, "validate patient failed")
	}

	patientBeforeUpdate, err := service.patientRepository.Get(patient.ID)
	if err != nil {
		log.Error().
			Err(err).
			Int("patient ID", patient.ID).
			Msg("get patient failed")

		return errors.Wrap(err, "get patient failed")
	}

	err = service.patientRepository.Update(patient)
	if err != nil {
		if isEntryNotValidErr(err) {
			return &modelValidationErr{err.Error()}
		}

		if isEntryNotExistErr(err) {
			return &modelNotExistErr{"patient doesn't exist"}
		}

		log.Error().
			Err(err).
			Msg("update patient failed")

		return errors.Wrap(err, "update patient failed")
	}

	// Patient status changed from another state to PatientStateCall
	if patient.Status == model.PatientStateCall && patient.Status != patientBeforeUpdate.Status {
		log.Debug().
			Int("pager", patient.PagerID).
			Msg("pager gets called")

		client := easycall.NewClient(service.cfg.EasyCall.URL, service.cfg.EasyCall.User, service.cfg.EasyCall.Password)

		if err := client.Send(&easycall.SendOptions{
			Receiver: patient.PagerID,
			Message:  "",
		}); err != nil {
			log.Error().
				Err(err).
				Int("patient ID", patient.ID).
				Int("pager ID", patient.PagerID).
				Msg("call pager failed")

			return &externalServiceErr{"easycall pager call failed"}
		}

		patient.Status = model.PatientStateCalled
	}

	return nil
}

// Remove deletes an existing patient
func (service *DefaultPatientService) Remove(patient *model.Patient) error {
	if err := service.patientRepository.Remove(patient); err != nil {
		if isEntryNotExistErr(err) {
			return &modelNotExistErr{"patient doesn't exist"}
		}

		log.Error().
			Err(err).
			Msg("remove patient failed")

		return errors.Wrap(err, "remove patient failed")
	}

	return nil
}
