package service

import (
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ClientService interface
type ClientService interface {
	GetAll() ([]*model.Client, error)
	GetByName(string) (*model.Client, error)
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

// GetByName returns a client by it's name
func (service *DefaultClientService) GetByName(name string) (*model.Client, error) {
	client, err := service.clientRepository.GetByName(name)
	if err != nil {
		log.Error().
			Err(err).
			Str("client name", name).
			Msg("get client failed")
	}

	return client, errors.Wrapf(err, "get client failed", name)
}
