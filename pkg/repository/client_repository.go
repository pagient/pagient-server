package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

type clientRepository struct {
	sqlRepository
}

// NewClientRepository returns a new instance of a ClientRepository
func NewClientRepository(db *gorm.DB) service.ClientRepository {
	return &clientRepository{sqlRepository{db}}
}

// GetAll returns all configured clients
func (repo *clientRepository) GetAll(sess service.DB) ([]*model.Client, error) {
	session := sess.(*gorm.DB)

	var clients []*model.Client
	err := session.Find(&clients).Error

	return clients, errors.Wrap(err, "select all clients failed")
}

// Get returns a client by it's id
func (repo *clientRepository) Get(sess service.DB, id uint) (*model.Client, error) {
	session := sess.(*gorm.DB)

	client := &model.Client{}
	err := session.First(client, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return client, errors.Wrap(err, "select client by id failed")
}

func (repo *clientRepository) GetByUser(sess service.DB, username string) (*model.Client, error) {
	session := sess.(*gorm.DB)

	client := &model.Client{}
	err := session.
		Joins("JOIN users ON users.client_id = clients.id").
		Where("users.username = ?", username).First(client).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return client, errors.Wrap(err, "select client by user failed")
}
