package filedb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/satori/go.uuid"
)

const (
	patientCollection = "patient"
)

var (
	lock *sync.Mutex
)

type driver interface {
	Write(string, string, interface{}) error
	Read(string, string, interface{}) error
	ReadAll(string string) ([]string, error)
	Delete(string, string) error
}

// FileDatabase struct
type FileDatabase struct {
	driver driver
}

// GetPatient loads a patient by ID
func (db *FileDatabase) GetPatient(id uuid.UUID) (*model.Patient, error) {
	lock.Lock()
	defer lock.Unlock()

	patient := &model.Patient{}
	if err := db.driver.Read(patientCollection, id.String(), patient); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("patient not found")
		}
		return nil, err
	}

	return patient, nil
}

// GetPatients loads all patients
func (db *FileDatabase) GetPatients() ([]*model.Patient, error) {
	lock.Lock()
	defer lock.Unlock()

	records, err := db.driver.ReadAll(patientCollection)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	patients := []*model.Patient{}
	for _, p := range records {
		patient := model.Patient{}
		if err := json.Unmarshal([]byte(p), &patient); err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}

// AddPatient persists a patient
func (db *FileDatabase) AddPatient(patient *model.Patient) error {
	lock.Lock()
	defer lock.Unlock()

	patient.ID = uuid.NewV4()
	if err := db.driver.Write(patientCollection, patient.ID.String(), patient); err != nil {
		patient.ID = uuid.Nil
		return err
	}

	return nil
}

// UpdatePatient updates a persistent patient
func (db *FileDatabase) UpdatePatient(patient *model.Patient) error {
	if patient.ID == uuid.Nil {
		return fmt.Errorf("Failure trying to update a patient without ID")
	}

	lock.Lock()
	defer lock.Unlock()

	if err := db.driver.Delete(patientCollection, patient.ID.String()); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("patient not found")
		}
		return err
	}

	err := db.driver.Write(patientCollection, patient.ID.String(), patient)

	return err
}

// RemovePatient removes a persistent patient
func (db *FileDatabase) RemovePatient(patient *model.Patient) error {
	lock.Lock()
	defer lock.Unlock()

	err := db.driver.Delete(patientCollection, patient.ID.String())
	if os.IsNotExist(err) {
		return fmt.Errorf("patient not found")
	}

	return err
}

// New creates and returns a new file database connection
func New(rootPath string) (*FileDatabase, error) {
	fileDatabase := new(FileDatabase)

	// Set up scribble json file store
	db, err := scribble.New(rootPath, nil)
	if err != nil {
		return nil, err
	}
	fileDatabase.driver = db

	return fileDatabase, nil
}

func init() {
	lock = &sync.Mutex{}
}
