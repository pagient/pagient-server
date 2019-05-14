package model

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

// Client struct
type Client struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null;unique"`
}

// Validate validates the client
func (client *Client) Validate() error {
	if err := validation.ValidateStruct(client,
		validation.Field(&client.Name, validation.Required, validation.Match(regexp.MustCompile("[[:print:]]+$"))),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occurred")
		}

		return &modelValidationErr{err.Error()}
	}

	return nil
}
