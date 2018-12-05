package service

import (
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetByUser returns all active tokens by username
func (service *DefaultService) ListTokensByUser(username string) ([]*model.Token, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	tokens, err := tx.GetTokenByUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Str("user", username).
			Msg("get token failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get token by user failed")
	}

	tx.Commit()
	return tokens, nil
}

// Get returns a token
func (service *DefaultService) ShowToken(rawToken string) (*model.Token, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	token, err := tx.GetToken(rawToken)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get token failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get token failed")
	}

	tx.Commit()
	return token, nil
}

// Add adds an active token to a user
func (service *DefaultService) AddToken(token *model.Token) error {
	tx, err := service.db.Begin()
	if err != nil {
		return errors.Wrap(err, "create transaction failed")
	}

	token, err = tx.AddToken(token)
	if err != nil {
		log.Error().
			Err(err).
			Msg("add token failed")

		tx.Rollback()
		return errors.Wrap(err, "add token failed")
	}

	tx.Commit()
	return nil
}

// Remove removes an active token from a user
func (service *DefaultService) DeleteToken(token *model.Token) error {
	tx, err := service.db.Begin()
	if err != nil {
		return errors.Wrap(err, "create transaction failed")
	}

	err = tx.RemoveToken(token)
	if err != nil {
		tx.Rollback()

		if isEntryNotExistErr(err) {
			return &modelNotExistErr{"token doesn't exist"}
		}

		log.Error().
			Err(err).
			Msg("remove token failed")

		return errors.Wrap(err, "remove token failed")
	}

	tx.Commit()
	return nil
}
