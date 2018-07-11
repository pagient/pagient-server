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
	clients := make([]*model.Client, len(repo.cfg.General.Clients))
	for i, clientInfo := range repo.cfg.General.Clients {
		pair := strings.SplitN(clientInfo, ":", 2)

		id, err := strconv.Atoi(pair[0])
		if err != nil {
			return nil, errors.Wrap(err, "integer string conversion failed")
		}

		clients[i] = &model.Client{ID: id, Name: pair[1]}
	}

	return clients, nil
}

// Get returns a client by it's id
func (repo *clientCfgRepository) Get(id int) (*model.Client, error) {
	clients, err := repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "get clients failed")
	}

	for _, client := range clients {
		if client.ID == id {
			return client, nil
		}
	}

	return nil, nil
}

func (repo *clientCfgRepository) GetByUser(user *model.User) (*model.Client, error) {
	for _, userClientInfo := range repo.cfg.General.UserClient {
		pair := strings.SplitN(userClientInfo, ":", 2)

		if pair[0] == user.Username {
			id, err := strconv.Atoi(pair[1])
			if err != nil {
				return nil, errors.Wrap(err, "integer string conversion failed")
			}

			return repo.Get(id)
		}
	}

	return nil, nil
}
