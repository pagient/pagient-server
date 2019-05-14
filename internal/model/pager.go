package model

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

// Pager struct
type Pager struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"not null;unique"`
	EasyCallID uint   `gorm:"not null;unique"`
}

// Validate validates the pager
func (pager *Pager) Validate() error {
	if err := validation.ValidateStruct(pager,
		validation.Field(&pager.Name, validation.Required, validation.Match(regexp.MustCompile("[[:print:]]+$"))),
		validation.Field(&pager.EasyCallID, validation.Required),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occurred")
		}

		return &modelValidationErr{err.Error()}
	}

	return nil
}
