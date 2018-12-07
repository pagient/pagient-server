package service

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ListClients returns all clients
func (service *defaultService) ListClients() ([]*model.Client, error) {
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

// ShowClient returns a client by it's id
func (service *defaultService) ShowClient(id uint) (*model.Client, error) {
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

// ShowClientByUser returns a client belonging to the given user
func (service *defaultService) ShowClientByUser(username string) (*model.Client, error) {
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

// CreateClient creates a new client
func (service *defaultService) CreateClient(client *model.Client) (*model.Client, error) {
	if err := service.validateClient(client); err != nil {
		return nil, errors.WithStack(err)
	}

	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	client, err = tx.AddClient(client)
	if err != nil {
		log.Error().
			Err(err).
			Msg("add client failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "add client failed")
	}

	tx.Commit()
	return client, nil
}

func (service *defaultService) validateClient(client *model.Client) error {
	if err := client.Validate(); err != nil {
		if model.IsValidationErr(err) {
			return &modelValidationErr{err.Error()}
		}

		return errors.Wrap(err, "validate client failed")
	}

	return nil
}
