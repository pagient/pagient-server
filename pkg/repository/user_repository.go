package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/pagient/pagient-server/pkg/service"
)

type userRepository struct {
	sqlRepository
}

// NewUserRepository returns a new instance of a UserRepository
func NewUserRepository(db *gorm.DB) service.UserRepository {
	return &userRepository{sqlRepository{db}}
}

func (repo *userRepository) GetAll(sess service.DB) ([]*model.User, error) {
	session := sess.(*gorm.DB)

	var users []*model.User
	err := session.Find(&users).Error

	return users, errors.Wrap(err, "select all users failed")
}

func (repo *userRepository) Get(sess service.DB, username string) (*model.User, error) {
	session := sess.(*gorm.DB)

	user := &model.User{}
	err := session.Where(&model.User{
		Username: username,
	}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by username failed")
}

func (repo *userRepository) GetByToken(sess service.DB, rawToken string) (*model.User, error) {
	session := sess.(*gorm.DB)

	user := &model.User{}
	err := session.
		Joins("JOIN tokens ON tokens.user_id = users.id").
		Where("tokens.raw = ?", rawToken).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by token failed")
}
