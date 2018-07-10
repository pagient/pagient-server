package service

import (
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// TokenService interface
type TokenService interface {
	Get(string) (*model.Token, error)
	Add(string, *model.Token) error
	Remove(string) error
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

func (service *DefaultTokenService) Get(username string) (*model.Token, error) {
	token, err := service.tokenRepository.Get(username)
	if err != nil {
		log.Error().
			Err(err).
			Str("user", username).
			Msg("get token failed")
	}

	return token, errors.Wrap(err, "get token failed")
}

func (service *DefaultTokenService) Add(username string, token *model.Token) error {
	err := service.tokenRepository.Add(username, token)
	if err != nil {
		if isEntryExistErr(err) {
			return &modelExistErr{"token already exists"}
		}

		log.Error().
			Err(err).
			Msg("add token failed")
	}

	return errors.Wrap(err, "add token failed")
}

func (service *DefaultTokenService) Remove(username string) error {
	if err := service.tokenRepository.Remove(username); err != nil {
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
