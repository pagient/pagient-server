package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type patientRepository struct {
	sqlRepository
}

// NewPatientRepository returns a new instance of a PatientRepository
func NewPatientRepository(db *gorm.DB) service.PatientRepository {
	return &patientRepository{sqlRepository{db}}
}

// GetAll lists all patients
func (repo *patientRepository) GetAll(sess service.DB) ([]*model.Patient, error) {
	session := sess.(*gorm.DB)

	var patients []*model.Patient
	err := session.Find(&patients).Error

	return patients, errors.Wrap(err, "select all patients failed")
}

// Get returns a patient by ID
func (repo *patientRepository) Get(sess service.DB, id uint) (*model.Patient, error) {
	session := sess.(*gorm.DB)

	patient := &model.Patient{}
	err := session.Find(patient, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return patient, errors.Wrap(err, "select patient by id failed")
}

// Add stores the values in the repository
func (repo *patientRepository) Add(sess service.DB, patient *model.Patient) (*model.Patient, error) {
	session := sess.(*gorm.DB)

	// FIXME: handle sql constraint errors
	err := session.Create(patient).Error

	return patient, errors.Wrap(err, "create patient failed")
}

// Update updates the values in the repository
func (repo *patientRepository) Update(sess service.DB, patient *model.Patient) (*model.Patient, error) {
	session := sess.(*gorm.DB)

	// FIXME: handle sql constraint errors
	err := session.Save(patient).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, &entryNotExistErr{"patient not found"}
	}

	return patient, errors.Wrap(err, "update patient failed")
}

// Remove deletes the values from the repository
func (repo *patientRepository) Remove(sess service.DB, patient *model.Patient) (*model.Patient, error) {
	session := sess.(*gorm.DB)

	err := session.Delete(patient).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, &entryNotExistErr{"patient not found"}
	}

	return patient, errors.Wrap(err, "delete patient failed")
}

// MarkAllExceptPatientInactiveByPatientClient sets active to false for every patient by that client
func (repo *patientRepository) MarkAllExceptPatientInactiveByPatientClient(sess service.DB, patient *model.Patient) ([]*model.Patient, error) {
	session := sess.(*gorm.DB)

	statement := session.Where(&model.Patient{
		ClientID: patient.ClientID,
		Active:   true,
	}).Not(&model.Patient{
		ID: patient.ID,
	})

	var patients []*model.Patient
	if err := statement.Find(&patients).Error; err != nil {
		log.Error().
			Err(err).
			Msg("find active patients by client failed")
	}

	if err := statement.Model(model.Patient{}).Updates(map[string]interface{}{"active": false}).Error; err != nil {
		log.Error().
			Err(err).
			Msg("update active patients by client failed")
	}

	return patients, nil
}

// RemoveAllExceptPatientInactiveNoPagerByPatientClient deletes the patients that are inactive, have no pager assigned and are from that client
func (repo *patientRepository) RemoveAllExceptPatientInactiveNoPagerByPatientClient(sess service.DB, patient *model.Patient) ([]*model.Patient, error) {
	session := sess.(*gorm.DB)

	statement := session.Where(&model.Patient{
		ClientID: patient.ClientID,
		Active:   false,
	}).Where("pager_id = 0").Not(&model.Patient{
		ID: patient.ID,
	})

	var patients []*model.Patient
	if err := statement.Find(&patients).Error; err != nil {
		log.Error().
			Err(err).
			Msg("find active patients w/o pager by client failed")
	}

	if err := statement.Delete(model.Patient{}).Error; err != nil {
		log.Error().
			Err(err).
			Msg("delete active patients w/o pager by client failed")
	}

	return patients, nil
}
