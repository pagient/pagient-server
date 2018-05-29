package model

import (
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/satori/go.uuid"
)

var cfg *config.Config
var db Database

// Database interface models adhere to
type Database interface {
	GetPatient(id uuid.UUID) (*Patient, error)
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
