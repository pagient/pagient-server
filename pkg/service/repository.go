package service

import (
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
)

// DB generic interface
type DB interface{}

// ClientRepository interface
type ClientRepository interface {
	GetAll() ([]*model.Client, error)
	Get(uint) (*model.Client, error)
	GetByUser(string) (*model.Client, error)
}

// PagerRepository interface
type PagerRepository interface {
	GetAll() ([]*model.Pager, error)
	GetUnassigned() ([]*model.Pager, error)
	Get(uint) (*model.Pager, error)
}

// PatientRepository interface
type PatientRepository interface {
	BeginTx() DB
	RollbackTx(DB) DB
	CommitTx(DB) DB
	GetAll(DB) ([]*model.Patient, error)
	Get(DB, uint) (*model.Patient, error)
	Add(DB, *model.Patient) (*model.Patient, error)
	Update(DB, *model.Patient) (*model.Patient, error)
	Remove(DB, *model.Patient) (*model.Patient, error)
	MarkAllExceptPatientInactiveByPatientClient(DB, *model.Patient) ([]*model.Patient, error)
	RemoveAllExceptPatientInactiveNoPagerByPatientClient(DB, *model.Patient) ([]*model.Patient, error)
}

// TokenRepository interface
type TokenRepository interface {
	Get(string) (*model.Token, error)
	GetByUser(string) ([]*model.Token, error)
	Add(*model.Token) (*model.Token, error)
	Remove(*model.Token) (*model.Token, error)
}

// UserRepository interface
type UserRepository interface {
	GetAll() ([]*model.User, error)
	Get(string) (*model.User, error)
	GetByToken(string) (*model.User, error)
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
