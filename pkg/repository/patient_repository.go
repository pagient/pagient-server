package repository

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/pkg/errors"
)

const (
	patientCollection = "patient"
)

var (
	lock *sync.Mutex

	patientRepositoryOnce     sync.Once
	patientRepositoryInstance service.PatientRepository
)

type driver interface {
	Write(string, string, interface{}) error
	Read(string, string, interface{}) error
	ReadAll(string string) ([]string, error)
	Delete(string, string) error
}

// GetPatientRepositoryInstance creates and returns a new PatientFileRepository
func GetPatientRepositoryInstance(cfg *config.Config) (service.PatientRepository, error) {
	var err error

	patientRepositoryOnce.Do(func() {
		// Set up scribble json file store
		var db driver
		db, err = scribble.New(cfg.General.Root, nil)

		lock = &sync.Mutex{}

		patientRepositoryInstance = &patientFileRepository{
			db: db,
		}
	})

	if err != nil {
		return nil, errors.Wrap(err, "init scribble store failed")
	}

	return patientRepositoryInstance, nil
}

type patientFileRepository struct {
	db driver
}

// GetAll lists all patients
func (repo *patientFileRepository) GetAll() ([]*model.Patient, error) {
	lock.Lock()
	defer lock.Unlock()

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
	lock.Lock()
	defer lock.Unlock()

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
func (repo *patientFileRepository) Add(patient *model.Patient) error {
	lock.Lock()
	defer lock.Unlock()

	if patient.ID == 0 {
		return &entryNotValidErr{"id: cannot be blank"}
	}

	pat := &model.Patient{}
	if err := repo.db.Read(patientCollection, strconv.Itoa(patient.ID), pat); err != nil && !isNotFoundErr(err) {
		return errors.Wrap(err, "read patient failed")
	}
	if pat.ID != 0 {
		return &entryExistErr{"patient already exists"}
	}

	err := repo.db.Write(patientCollection, strconv.Itoa(patient.ID), patient)
	return errors.Wrap(err, "write patient failed")
}

// Update updates the values in the repository
func (repo *patientFileRepository) Update(patient *model.Patient) error {
	lock.Lock()
	defer lock.Unlock()

	if patient.ID == 0 {
		return &entryNotValidErr{"id: cannot be blank"}
	}

	if err := repo.db.Read(patientCollection, strconv.Itoa(patient.ID), &model.Patient{}); err != nil {
		if isNotFoundErr(err) {
			return &entryNotExistErr{"patient not found"}
		}
		return errors.Wrap(err, "read patient failed")
	}

	err := repo.db.Write(patientCollection, strconv.Itoa(patient.ID), patient)
	return errors.Wrap(err, "write patient failed")
}

// Remove deletes the values from the repository
func (repo *patientFileRepository) Remove(patient *model.Patient) error {
	lock.Lock()
	defer lock.Unlock()

	err := repo.db.Delete(patientCollection, strconv.Itoa(patient.ID))
	if err != nil {
		if isNotFoundErr(err) {
			return &entryNotExistErr{"patient not found"}
		}
		return errors.Wrap(err, "delete patient failed")
	}
	return nil
}
