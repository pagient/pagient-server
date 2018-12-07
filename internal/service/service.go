package service

import (
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/notifier"
)

// ClientService interface
type ClientService interface {
	ListClients() ([]*model.Client, error)
	ShowClient(uint) (*model.Client, error)
	ShowClientByUser(string) (*model.Client, error)
	CreateClient(*model.Client) (*model.Client, error)
}

// PagerService interface
type PagerService interface {
	ListPagers() ([]*model.Pager, error)
	ShowPager(uint) (*model.Pager, error)
}

// PatientService interface
type PatientService interface {
	ListPatients() ([]*model.Patient, error)
	ListPagerPatientsByStatus(...model.PatientState) ([]*model.Patient, error)
	ShowPatient(uint) (*model.Patient, error)
	CreatePatient(*model.Patient) (*model.Patient, error)
	UpdatePatient(*model.Patient) (*model.Patient, error)
	DeletePatient(*model.Patient) error
	CallPatient(*model.Patient) error
}

// TokenService interface
type TokenService interface {
	ListTokensByUser(string) ([]*model.Token, error)
	ShowToken(string) (*model.Token, error)
	CreateToken(*model.Token) error
	DeleteToken(*model.Token) error
}

// UserService interface
type UserService interface {
	ListUsers() ([]*model.User, error)
	ShowUser(string) (*model.User, error)
	ShowUserByToken(string) (*model.User, error)
	CreateUser(*model.User) (*model.User, error)
	ChangeUserPassword(*model.User) (*model.User, error)
	Login(string, string) (*model.User, bool, error)
}

// Service interface combines all concrete model services
type Service interface {
	ClientService
	PagerService
	PatientService
	TokenService
	UserService
}

type defaultService struct {
	db       DB
	notifier notifier.Notifier
}

// NewService constructs a new service layer
func NewService(db DB, notifier notifier.Notifier) Service {
	return &defaultService{db, notifier}
}
