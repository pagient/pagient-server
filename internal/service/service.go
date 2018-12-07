package service

import (
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/notifier"
)

type ClientService interface {
	ListClients() ([]*model.Client, error)
	ShowClient(uint) (*model.Client, error)
	ShowClientByUser(string) (*model.Client, error)
	CreateClient(*model.Client) (*model.Client, error)
}

type PagerService interface {
	ListPagers() ([]*model.Pager, error)
	ShowPager(uint) (*model.Pager, error)
}

type PatientService interface {
	ListPatients() ([]*model.Patient, error)
	ListPagerPatientsByStatus(...model.PatientState) ([]*model.Patient, error)
	ShowPatient(uint) (*model.Patient, error)
	AddPatient(*model.Patient) (*model.Patient, error)
	UpdatePatient(*model.Patient) (*model.Patient, error)
	DeletePatient(*model.Patient) error
	CallPatient(*model.Patient) error
}

type TokenService interface {
	ListTokensByUser(string) ([]*model.Token, error)
	ShowToken(string) (*model.Token, error)
	AddToken(*model.Token) error
	DeleteToken(*model.Token) error
}

type UserService interface {
	ListUsers() ([]*model.User, error)
	ShowUser(string) (*model.User, error)
	ShowUserByToken(string) (*model.User, error)
	CreateUser(*model.User) (*model.User, error)
	ChangeUserPassword(*model.User) (*model.User, error)
	Login(string, string) (*model.User, bool, error)
}

type Service interface {
	ClientService
	PagerService
	PatientService
	TokenService
	UserService
}

type DefaultService struct {
	db       DB
	notifier notifier.Notifier
}

func Init(db DB, notifier notifier.Notifier) *DefaultService {
	return &DefaultService{db, notifier}
}
