package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
)

// User struct
type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Client   Client `gorm:"save_associations:false"`
	ClientID uint   `gorm:"unique"`
}

// Validate validates the user
func (user *User) Validate(clients []*Client) error {
	// convert pager slice to generic interface slice
	clientIDs := make([]interface{}, len(clients))
	for i, client := range clients {
		clientIDs[i] = client.ID
	}

	if err := validation.ValidateStruct(user,
		validation.Field(&user.Username, validation.Required, is.Alphanumeric),
		validation.Field(&user.Password, validation.Required, validation.Length(1, 100)),
		validation.Field(&user.ClientID, validation.In(clientIDs...)),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occured")
		}

		return &modelValidationErr{err.Error()}
	}

	return nil
}

func (user *User) ValidatePasswordChange() error {
	if err := validation.ValidateStruct(user,
		validation.Field(&user.ID, validation.Required),
		validation.Field(&user.Password, validation.Required, validation.Length(1, 100)),
	); err != nil {
		if e, ok := err.(validation.InternalError); ok {
			return errors.Wrap(e, "internal validation error occured")
		}

		return &modelValidationErr{err.Error()}
	}
}
