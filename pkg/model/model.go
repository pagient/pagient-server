package model

var db Storage

type Storage interface {
	AddPatient(patient *Patient) error
	GetPatient(id int64) (*Patient, error)
	GetPatients() ([]*Patient, error)
	RemovePatient(patient *Patient) error
}

func Init(database Storage) {
	db = database
}
