package service

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetAll returns all clients
func (service *DefaultService) ListClients() ([]*model.Client, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	clients, err := tx.GetClients()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all clients failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get all clients failed")
	}

	tx.Commit()
	return clients, nil
}

// Get returns a client by it's id
func (service *DefaultService) ShowClient(id uint) (*model.Client, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	client, err := tx.GetClient(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("client id", id).
			Msg("get client failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get client failed")
	}

	tx.Commit()
	return client, nil
}

// GetByUser returns a client belonging to the given user
func (service *DefaultService) ShowClientByUser(username string) (*model.Client, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	client, err := tx.GetClientByUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", username).
			Msg("get client by user failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get client by user failed")
	}

	tx.Commit()
	return client, nil
}
