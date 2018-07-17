package repository

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

var (
	clientRepositoryOnce     sync.Once
	clientRepositoryInstance service.ClientRepository
)

// GetClientRepositoryInstance creates and returns a new ClientCfgRepository
func GetClientRepositoryInstance(db *gorm.DB) (service.ClientRepository, error) {
	clientRepositoryOnce.Do(func() {
		clientRepositoryInstance = &clientRepository{db}
	})

	return clientRepositoryInstance, nil
}

type clientRepository struct {
	db *gorm.DB
}

// GetAll returns all configured clients
func (repo *clientRepository) GetAll() ([]*model.Client, error) {
	var clients []*model.Client
	err := repo.db.Find(&clients).Error

	return clients, errors.Wrap(err, "select all clients failed")
}

// Get returns a client by it's id
func (repo *clientRepository) Get(id uint) (*model.Client, error) {
	client := &model.Client{}
	err := repo.db.First(client, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return client, errors.Wrap(err, "select client by id failed")
}

func (repo *clientRepository) GetByUser(username string) (*model.Client, error) {
	client := &model.Client{}
	err := repo.db.
		Joins("JOIN users ON users.client_id = clients.id").
		Where("users.username = ?", username).First(client).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return client, errors.Wrap(err, "select client by user failed")
}
