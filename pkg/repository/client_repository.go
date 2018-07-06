package repository

import (
	"strconv"
	"strings"
	"sync"

	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/pkg/errors"
)

var (
	clientRepositoryOnce     sync.Once
	clientRepositoryInstance service.ClientRepository
)

// GetClientRepositoryInstance creates and returns a new ClientCfgRepository
func GetClientRepositoryInstance(cfg *config.Config) (service.ClientRepository, error) {
	clientRepositoryOnce.Do(func() {
		clientRepositoryInstance = &clientCfgRepository{cfg}
	})

	return clientRepositoryInstance, nil
}

type clientCfgRepository struct {
	cfg *config.Config
}

// GetAll returns all configured clients
func (repo *clientCfgRepository) GetAll() ([]*model.Client, error) {
	var clients []*model.Client
	for _, clientInfo := range repo.cfg.General.Clients {
		pair := strings.SplitN(clientInfo, ":", 2)

		id, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, errors.Wrap(err, "integer string conversion failed")
		}

		clients = append(clients, &model.Client{ID: id, Name: pair[0]})
	}

	return clients, nil
}

// GetByName returns a client by name
func (repo *clientCfgRepository) GetByName(name string) (*model.Client, error) {
	clientID, err := repo.getClientID(name)
	if err != nil {
		return nil, errors.Wrap(err, "get client id failed")
	}

	if clientID == 0 {
		return nil, nil
	}

	return &model.Client{
		ID:   clientID,
		Name: name,
	}, nil
}

func (repo *clientCfgRepository) getClientID(name string) (int, error) {
	for _, clientMapping := range repo.cfg.General.Clients {
		clientInfo := strings.SplitN(clientMapping, ":", 2)
		if clientInfo[0] == name {
			id, err := strconv.Atoi(clientInfo[1])
			return id, errors.Wrap(err, "integer string conversion failed")
		}
	}

	return 0, nil
}
