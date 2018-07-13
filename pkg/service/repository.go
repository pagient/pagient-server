package service

import (
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
)

// ClientRepository interface
type ClientRepository interface {
	GetAll() ([]*model.Client, error)
	Get(int) (*model.Client, error)
	GetByUser(*model.User) (*model.Client, error)
}

// PagerRepository interface
type PagerRepository interface {
	GetAll() ([]*model.Pager, error)
	Get(int) (*model.Pager, error)
}

// PatientRepository interface
type PatientRepository interface {
	GetAll() ([]*model.Patient, error)
	Get(int) (*model.Patient, error)
	Add(*model.Patient) error
	Update(*model.Patient) error
	Remove(*model.Patient) error
	MarkAllExceptPatientInactiveByPatientClient(*model.Patient) error
	RemoveAllExceptPatientInactiveNoPagerByPatientClient(*model.Patient) error
}

// TokenRepository interface
type TokenRepository interface {
	Get(string) ([]*model.Token, error)
	Add(*model.Token) error
	Remove(*model.Token) error
}

// UserRepository interface
type UserRepository interface {
	GetAll() ([]*model.User, error)
	Get(string) (*model.User, error)
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
