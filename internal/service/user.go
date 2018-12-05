package service

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetAll returns all users
func (service *DefaultService) ListUsers() ([]*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	users, err := tx.GetUsers()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all users failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get all users failed")
	}

	tx.Commit()
	return users, nil
}

// Get returns a user by it's username
func (service *DefaultService) ShowUser(username string) (*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	user, err := tx.GetUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get user failed")
	}

	tx.Commit()
	return user, nil
}

// GetByToken returns a user by token
func (service *DefaultService) ShowUserByToken(rawToken string) (*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	user, err := tx.GetUserByToken(rawToken)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user by token failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get user by token failed")
	}

	tx.Commit()
	return user, nil
}

// Login checks whether the combination of username and password is valid
func (service *DefaultService) Login(username, password string) (*model.User, bool, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, false, errors.Wrap(err, "create transaction failed")
	}

	user, err := tx.GetUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		tx.Rollback()
		return nil, false, errors.Wrap(err, "get user failed")
	}

	tx.Commit()
	return user, user != nil && user.Password == password, nil
}
