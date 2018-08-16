package service

import (
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
)

// DB generic interface
type DB interface{}

type repository interface {
	BeginTx() DB
	RollbackTx(DB) DB
	CommitTx(DB) DB
}

// ClientRepository interface
type ClientRepository interface {
	repository
	GetAll(DB) ([]*model.Client, error)
	Get(DB, uint) (*model.Client, error)
	GetByUser(DB, string) (*model.Client, error)
}

// PagerRepository interface
type PagerRepository interface {
	repository
	GetAll(DB) ([]*model.Pager, error)
	GetUnassigned(DB) ([]*model.Pager, error)
	Get(DB, uint) (*model.Pager, error)
}

// PatientRepository interface
type PatientRepository interface {
	repository
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
	repository
	Get(DB, string) (*model.Token, error)
	GetByUser(DB, string) ([]*model.Token, error)
	Add(DB, *model.Token) (*model.Token, error)
	Remove(DB, *model.Token) (*model.Token, error)
}

// UserRepository interface
type UserRepository interface {
	repository
	GetAll(DB) ([]*model.User, error)
	Get(DB, string) (*model.User, error)
	GetByToken(DB, string) (*model.User, error)
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
