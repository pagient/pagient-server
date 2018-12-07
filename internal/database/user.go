package database

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pkg/errors"
)

// GetUsers returns all users
func (t *tx) GetUsers() ([]*model.User, error) {
	var users []*model.User
	err := t.Find(&users).Error

	return users, errors.Wrap(err, "select all users failed")
}

// GetUser returns a user by username
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

// GetUserByToken returns a user with given token
func (t *tx) GetUserByToken(rawToken string) (*model.User, error) {
	user := &model.User{}
	err := t.Joins("JOIN tokens ON tokens.user_id = users.id").
		Where("tokens.raw = ?", rawToken).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by token failed")
}

// AddUser creates a new user
func (t *tx) AddUser(user *model.User) (*model.User, error) {
	// FIXME: handle sql constraint errors
	err := t.Create(user).Error

	return user, errors.Wrap(err, "create user failed")
}

// UpdateUserPassword updates only the password of provided user
func (t *tx) UpdateUserPassword(user *model.User) (*model.User, error) {
	// FIXME: handle sql constraint errors
	err := t.Model(user).UpdateColumn("password", user.Password).Error

	return user, errors.Wrap(err, "update password failed")
}
