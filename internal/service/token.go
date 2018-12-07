package service

import (
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ListTokensByUser returns all active tokens by username
func (service *defaultService) ListTokensByUser(username string) ([]*model.Token, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	tokens, err := tx.GetTokensByUser(username)
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

// ShowToken returns a token
func (service *defaultService) ShowToken(rawToken string) (*model.Token, error) {
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

// CreateToken adds an active token to a user
func (service *defaultService) CreateToken(token *model.Token) error {
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

// DeleteToken removes an active token from a user
func (service *defaultService) DeleteToken(token *model.Token) error {
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
