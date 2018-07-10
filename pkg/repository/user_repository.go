package repository

import (
	"strings"
	"sync"

	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/pkg/errors"
)

var (
	userRepositoryOnce     sync.Once
	userRepositoryInstance service.UserRepository
)

// GetUserRepositoryInstance creates and returns a new UserCfgRepository
func GetUserRepositoryInstance(cfg *config.Config) (service.UserRepository, error) {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = &userCfgRepository{cfg}
	})

	return userRepositoryInstance, nil
}

type userCfgRepository struct {
	cfg *config.Config
}

func (repo *userCfgRepository) GetAll() ([]*model.User, error) {
	users := make([]*model.User, len(repo.cfg.General.Users))
	for i, userInfo := range repo.cfg.General.Users {
		pair := strings.SplitN(userInfo, ":", 2)

		users[i] = &model.User{Username: pair[0], Password: pair[1]}
	}

	return users, nil
}

func (repo *userCfgRepository) Get(username string) (*model.User, error) {
	users, err := repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "get users failed")
	}

	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, nil
}
