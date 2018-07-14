package repository

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

const (
	patientCollection = "patient"
)

var (
	patientRepositoryOnce     sync.Once
	patientRepositoryInstance service.PatientRepository
)

// GetPatientRepositoryInstance creates and returns a new PatientFileRepository
func GetPatientRepositoryInstance(cfg *config.Config) (service.PatientRepository, error) {
	var err error

	patientRepositoryOnce.Do(func() {
		// Set up scribble json file store
		var db fileDriver
		db, err = scribble.New(cfg.General.Root, nil)

		patientRepositoryInstance = &patientFileRepository{
			lock: &sync.Mutex{},
			db:   db,
		}
	})

	if err != nil {
		return nil, errors.Wrap(err, "init scribble store failed")
	}

	return patientRepositoryInstance, nil
}

type patientFileRepository struct {
	lock *sync.Mutex
	db   fileDriver
}

// GetAll lists all patients
func (repo *patientFileRepository) GetAll() ([]*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	records, err := repo.db.ReadAll(patientCollection)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "read patients failed")
	}

	patients := make([]*model.Patient, len(records))
	for i, p := range records {
		patient := &model.Patient{}
		if err := json.Unmarshal([]byte(p), patient); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}
		patients[i] = patient
	}

	return patients, nil
}

// Get returns a patient by ID
func (repo *patientFileRepository) Get(id int) (*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	patient := &model.Patient{}
	if err := repo.db.Read(patientCollection, strconv.Itoa(id), patient); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "read patient failed")
	}

	return patient, nil
}

// Add stores the values in the repository
func (repo *patientFileRepository) Add(patient *model.Patient) (*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	if patient.ID == 0 {
		return nil, &entryNotValidErr{"id: cannot be blank"}
	}

	pat := &model.Patient{}
	if err := repo.db.Read(patientCollection, strconv.Itoa(patient.ID), pat); err != nil && !isNotFoundErr(err) {
		return nil, errors.Wrap(err, "read patient failed")
	}
	if pat.ID != 0 {
		return nil, &entryExistErr{"patient already exists"}
	}

	err := repo.db.Write(patientCollection, strconv.Itoa(patient.ID), patient)
	return patient, errors.Wrap(err, "write patient failed")
}

// Update updates the values in the repository
func (repo *patientFileRepository) Update(patient *model.Patient) (*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	if patient.ID == 0 {
		return nil, &entryNotValidErr{"id: cannot be blank"}
	}

	if err := repo.db.Read(patientCollection, strconv.Itoa(patient.ID), &model.Patient{}); err != nil {
		if isNotFoundErr(err) {
			return nil, &entryNotExistErr{"patient not found"}
		}
		return nil, errors.Wrap(err, "read patient failed")
	}

	err := repo.db.Write(patientCollection, strconv.Itoa(patient.ID), patient)
	return patient, errors.Wrap(err, "write patient failed")
}

// Remove deletes the values from the repository
func (repo *patientFileRepository) Remove(patient *model.Patient) (*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	err := repo.db.Delete(patientCollection, strconv.Itoa(patient.ID))
	if err != nil {
		if isNotFoundErr(err) {
			return nil, &entryNotExistErr{"patient not found"}
		}
		return nil, errors.Wrap(err, "delete patient failed")
	}

	return patient, nil
}

// MarkAllExceptPatientInactiveByPatientClient sets active to false for every patient by that client
func (repo *patientFileRepository) MarkAllExceptPatientInactiveByPatientClient(patient *model.Patient) ([]*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	records, err := repo.db.ReadAll(patientCollection)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "read patients failed")
	}

	patients := make([]*model.Patient, len(records))
	for i, p := range records {
		patient := &model.Patient{}
		if err := json.Unmarshal([]byte(p), patient); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}
		patients[i] = patient
	}

	var updatedPatients []*model.Patient
	for _, pat := range patients {
		if pat.ClientID == patient.ClientID && pat.ID != patient.ID && pat.Active {
			pat.Active = false

			if err := repo.db.Write(patientCollection, strconv.Itoa(pat.ID), pat); err != nil {
				return nil, errors.Wrap(err, "write patient failed")
			}

			updatedPatients = append(updatedPatients, pat)
		}
	}

	return updatedPatients, nil
}

// RemoveAllExceptPatientInactiveNoPagerByPatientClient deletes the patients that are inactive, have no pager assigned and are from that client
func (repo *patientFileRepository) RemoveAllExceptPatientInactiveNoPagerByPatientClient(patient *model.Patient) ([]*model.Patient, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	records, err := repo.db.ReadAll(patientCollection)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "read patients failed")
	}

	patients := make([]*model.Patient, len(records))
	for i, p := range records {
		patient := &model.Patient{}
		if err := json.Unmarshal([]byte(p), patient); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}
		patients[i] = patient
	}

	var deletedPatients []*model.Patient
	for _, pat := range patients {
		if pat.ClientID == patient.ClientID && pat.ID != patient.ID && !pat.Active && pat.PagerID == 0 {
			err := repo.db.Delete(patientCollection, strconv.Itoa(pat.ID))
			if err != nil && !isNotFoundErr(err) {
				return nil, errors.Wrap(err, "delete patient failed")
			}

			deletedPatients = append(deletedPatients, pat)
		}
	}

	return deletedPatients, nil
}
