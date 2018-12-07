package database

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetAll returns all configured clients
func (t *tx) GetClients() ([]*model.Client, error) {
	var clients []*model.Client
	err := t.Find(&clients).Error

	return clients, errors.Wrap(err, "select all clients failed")
}

// Get returns a client by it's id
func (t *tx) GetClient(id uint) (*model.Client, error) {
	client := &model.Client{}
	err := t.First(client, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return client, errors.Wrap(err, "select client by id failed")
}

func (t *tx) GetClientByUser(username string) (*model.Client, error) {
	client := &model.Client{}
	err := t.Joins("JOIN users ON users.client_id = clients.id").
		Where("users.username = ?", username).First(client).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return client, errors.Wrap(err, "select client by user failed")
}

func (t *tx) AddClient(client *model.Client) (*model.Client, error) {
	// FIXME: handle sql constraint errors
	err := t.Create(client).Error

	return client, errors.Wrap(err, "create client failed")
}
