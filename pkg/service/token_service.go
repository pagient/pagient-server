package service

import (
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// TokenService interface
type TokenService interface {
	Get(string) (*model.Token, error)
	GetByUser(string) ([]*model.Token, error)
	Add(*model.Token) error
	Remove(*model.Token) error
}

// DefaultTokenService struct
type DefaultTokenService struct {
	cfg             *config.Config
	tokenRepository TokenRepository
}

// NewTokenService initializes a TokenService
func NewTokenService(cfg *config.Config, tokenRepository TokenRepository) TokenService {
	return &DefaultTokenService{
		cfg:             cfg,
		tokenRepository: tokenRepository,
	}
}

// Get returns a token
func (service *DefaultTokenService) Get(rawToken string) (*model.Token, error) {
	token, err := service.tokenRepository.Get(rawToken)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get token failed")
	}

	return token, errors.Wrap(err, "get token failed")
}

// GetByUser returns all active tokens by username
func (service *DefaultTokenService) GetByUser(username string) ([]*model.Token, error) {
	tokens, err := service.tokenRepository.GetByUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Str("user", username).
			Msg("get token failed")
	}

	return tokens, errors.Wrap(err, "get token by user failed")
}

// Add adds an active token to a user
func (service *DefaultTokenService) Add(token *model.Token) error {
	token, err := service.tokenRepository.Add(token)
	if err != nil {
		log.Error().
			Err(err).
			Msg("add token failed")
	}

	return errors.Wrap(err, "add token failed")
}

// Remove removes an active token from a user
func (service *DefaultTokenService) Remove(token *model.Token) error {
	token, err := service.tokenRepository.Remove(token)
	if err != nil {
		if isEntryNotExistErr(err) {
			return &modelNotExistErr{"token doesn't exist"}
		}

		log.Error().
			Err(err).
			Msg("remove token failed")

		return errors.Wrap(err, "remove token failed")
	}

	return nil
}
