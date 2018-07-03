package model

import (
	"github.com/pagient/pagient-api/pkg/config"
)

var cfg *config.Config
var db Database

// Database interface models adhere to
type Database interface {
	GetPatient(id int) (*Patient, error)
	GetPatients() ([]*Patient, error)
	AddPatient(patient *Patient) error
	UpdatePatient(patient *Patient) error
	RemovePatient(patient *Patient) error
}

// Init initializes the models with the config and the configured database connection
func Init(config *config.Config, database Database) {
	cfg = config
	db = database
}
