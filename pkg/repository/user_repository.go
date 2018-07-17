package repository

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

var (
	userRepositoryOnce     sync.Once
	userRepositoryInstance service.UserRepository
)

// GetUserRepositoryInstance creates and returns a new UserCfgRepository
func GetUserRepositoryInstance(db *gorm.DB) (service.UserRepository, error) {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = &userRepository{db}
	})

	return userRepositoryInstance, nil
}

type userRepository struct {
	db *gorm.DB
}

func (repo *userRepository) GetAll() ([]*model.User, error) {
	var users []*model.User
	err := repo.db.Find(&users).Error

	return users, errors.Wrap(err, "select all users failed")
}

func (repo *userRepository) Get(username string) (*model.User, error) {
	user := &model.User{}
	err := repo.db.Where(&model.User{
		Username: username,
	}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by username failed")
}

func (repo *userRepository) GetByToken(rawToken string) (*model.User, error) {
	user := &model.User{}
	err := repo.db.
		Joins("JOIN tokens ON tokens.user_id = users.id").
		Where("tokens.raw = ?", rawToken).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return user, errors.Wrap(err, "select user by token failed")
}
