package model

import "github.com/pkg/errors"

type validationErr interface {
	Valid() bool
}

type modelValidationErr struct {
	msg string
}

func (err *modelValidationErr) Error() string {
	return err.msg
}

func (err *modelValidationErr) Valid() bool {
	return true
}

// IsValidationErr returns true if err is a validation error
func IsValidationErr(err error) bool {
	ve, ok := errors.Cause(err).(validationErr)
	return ok && ve.Valid()
}
