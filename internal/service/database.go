package service

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
)

// DB interface
type DB interface {
	Begin() (Tx, error)
}

// Tx interface
type Tx interface {
	Commit() error
	Rollback() error

	ClientTx
	PagerTx
	PatientTx
	TokenTx
	UserTx
}

// ClientTx interface
type ClientTx interface {
	GetClients() ([]*model.Client, error)
	GetClient(uint) (*model.Client, error)
	GetClientByUser(string) (*model.Client, error)
	AddClient(*model.Client) (*model.Client, error)
}

// PagerTx interface
type PagerTx interface {
	GetPagers() ([]*model.Pager, error)
	GetUnassignedPagers() ([]*model.Pager, error)
	GetPager(uint) (*model.Pager, error)
}

// PatientTx interface
type PatientTx interface {
	GetPatients() ([]*model.Patient, error)
	GetPatientsWithPagerByStatus(...model.PatientState) ([]*model.Patient, error)
	// Get Patients by Client, Activity (first in slice) and Assignment of a Pager (second in slice)
	GetPatientsByClient(uint, ...bool) ([]*model.Patient, error)
	GetPatient(uint) (*model.Patient, error)
	AddPatient(*model.Patient) (*model.Patient, error)
	UpdatePatient(*model.Patient) (*model.Patient, error)
	MarkPatientsInactiveByClient(uint) error
	RemovePatient(*model.Patient) error
	// Remove Patients by Client, Activity (first in slice) and Assignment of a Pager (second in slice)
	RemovePatientsByClient(uint, ...bool) error
}

// TokenTx interface
type TokenTx interface {
	GetToken(string) (*model.Token, error)
	GetTokensByUser(string) ([]*model.Token, error)
	AddToken(*model.Token) (*model.Token, error)
	RemoveToken(*model.Token) error
}

// UserTx interface
type UserTx interface {
	GetUsers() ([]*model.User, error)
	GetUser(string) (*model.User, error)
	GetUserByToken(string) (*model.User, error)
	AddUser(*model.User) (*model.User, error)
	UpdateUserPassword(*model.User) (*model.User, error)
}

type entryExistErr interface {
	EntryExist() bool
}

func isEntryExistErr(err error) bool {
	ee, ok := errors.Cause(err).(entryExistErr)
	return ok && ee.EntryExist()
}

type entryNotExistErr interface {
	EntryNotExist() bool
}

func isEntryNotExistErr(err error) bool {
	ne, ok := errors.Cause(err).(entryNotExistErr)
	return ok && ne.EntryNotExist()
}

type entryNotValidErr interface {
	NotValid() bool
}

func isEntryNotValidErr(err error) bool {
	nv, ok := errors.Cause(err).(entryNotValidErr)
	return ok && nv.NotValid()
}
