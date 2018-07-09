package service

import (
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pkg/errors"
)

// ClientRepository interface
type ClientRepository interface {
	GetAll() ([]*model.Client, error)
	GetByName(string) (*model.Client, error)
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
