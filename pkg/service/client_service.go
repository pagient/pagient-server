package service

import (
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ClientService interface
type ClientService interface {
	GetAll() ([]*model.Client, error)
	Get(uint) (*model.Client, error)
	GetByUser(string) (*model.Client, error)
}

// DefaultClientService struct
type DefaultClientService struct {
	clientRepository ClientRepository
}

// NewClientService initializes a ClientService
func NewClientService(repository ClientRepository) ClientService {
	return &DefaultClientService{
		clientRepository: repository,
	}
}

// GetAll returns all clients
func (service *DefaultClientService) GetAll() ([]*model.Client, error) {
	clients, err := service.clientRepository.GetAll()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all clients failed")
	}

	return clients, errors.Wrap(err, "get all clients failed")
}

// Get returns a client by it's id
func (service *DefaultClientService) Get(id uint) (*model.Client, error) {
	client, err := service.clientRepository.Get(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("client id", id).
			Msg("get client failed")
	}

	return client, errors.Wrap(err, "get client failed")
}

// GetByUser returns a client belonging to the given user
func (service *DefaultClientService) GetByUser(username string) (*model.Client, error) {
	client, err := service.clientRepository.GetByUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", username).
			Msg("get client by user failed")
	}

	return client, errors.Wrap(err, "get client by user failed")
}
