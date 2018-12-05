package service

import "github.com/pkg/errors"

type existErr interface {
	Exist() bool
}

type modelExistErr struct {
	msg string
}

func (err *modelExistErr) Error() string {
	return err.msg
}

func (err *modelExistErr) Exist() bool {
	return true
}

// IsModelExistErr returns true if model already exists
func IsModelExistErr(err error) bool {
	ne, ok := errors.Cause(err).(existErr)
	return ok && ne.Exist()
}

type notExistErr interface {
	NotExist() bool
}

type modelNotExistErr struct {
	msg string
}

func (err *modelNotExistErr) Error() string {
	return err.msg
}

func (err *modelNotExistErr) NotExist() bool {
	return true
}

// IsModelNotExistErr returns true if model does not exist
func IsModelNotExistErr(err error) bool {
	ne, ok := errors.Cause(err).(notExistErr)
	return ok && ne.NotExist()
}

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

// IsModelValidationErr returns true if model is invalid
func IsModelValidationErr(err error) bool {
	ve, ok := errors.Cause(err).(validationErr)
	return ok && ve.Valid()
}

type serviceErr interface {
	Service() bool
}

type externalServiceErr struct {
	msg string
}

func (err *externalServiceErr) Error() string {
	return err.msg
}

func (err *externalServiceErr) Service() bool {
	return true
}

// IsExternalServiceErr returns true if external service raises an error
func IsExternalServiceErr(err error) bool {
	es, ok := errors.Cause(err).(serviceErr)
	return ok && es.Service()
}

type invalidArgErr interface {
	InvalidArgument() bool
}

type invalidArgumentErr struct {
	msg string
}

func (err *invalidArgumentErr) Error() string {
	return err.msg
}

func (err *invalidArgumentErr) InvalidArgument() bool {
	return true
}

// IsInvalidArgumentErr returns true if a given argument is invalid
func IsInvalidArgumentErr(err error) bool {
	ia, ok := errors.Cause(err).(invalidArgErr)
	return ok && ia.InvalidArgument()
}
