package database

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pkg/errors"
)

func (t *tx) GetUsers() ([]*model.User, error) {
	var users []*model.User
	err := t.Find(&users).Error

	return users, errors.Wrap(err, "select all users failed")
}

func (t *tx) GetUser(username string) (*model.User, error) {
	user := &model.User{}
	err := t.Where(&model.User{
		Username: username,
	}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by username failed")
}

func (t *tx) GetUserByToken(rawToken string) (*model.User, error) {
	user := &model.User{}
	err := t.Joins("JOIN tokens ON tokens.user_id = users.id").
		Where("tokens.raw = ?", rawToken).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by token failed")
}

func (t *tx) AddUser(user *model.User) (*model.User, error) {
	// FIXME: handle sql constraint errors
	err := t.Create(user).Error

	return user, errors.Wrap(err, "create user failed")
}
