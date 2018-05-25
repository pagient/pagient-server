package storage

import (
	"github.com/pagient/pagient-api/pkg/model"
)

type FileStorage struct {
	rootPath string
}

func (storage *FileStorage) AddPatient(patient *model.Patient) error {
	// TODO: implement

	return nil
}

func (storage *FileStorage) GetPatient(id int64) (*model.Patient, error) {
	// TODO: implement

	return nil, nil
}

func (storage *FileStorage) GetPatients() ([]*model.Patient, error) {
	// TODO: implement

	return nil, nil
}

func (storage *FileStorage) RemovePatient(patient *model.Patient) error {
	// TODO: implement

	return nil
}

func NewFileStorage(rootPath string) *FileStorage {
	fileStorage := new(FileStorage)
	fileStorage.rootPath = rootPath
	return fileStorage
}
