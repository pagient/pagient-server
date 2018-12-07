package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

// Client struct
type Client struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null;unique"`
}

// Validate validates the user
func (client *Client) Validate() error {
	if err := validation.ValidateStruct(client,
		validation.Field(&client.Name, validation.Required, is.Alphanumeric),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occured")
		}

		return &modelValidationErr{err.Error()}
	}

	return nil
}
