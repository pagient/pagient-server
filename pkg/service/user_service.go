package service

import (
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// UserService interface
type UserService interface {
	GetAll() ([]*model.User, error)
	Get(string) (*model.User, error)
	GetByToken(string) (*model.User, error)
	Login(username, password string) (*model.User, bool, error)
}

// DefaultUserService struct
type DefaultUserService struct {
	cfg            *config.Config
	userRepository UserRepository
}

// NewUserService initializes a PatientService
func NewUserService(cfg *config.Config, userRepository UserRepository) UserService {
	return &DefaultUserService{
		cfg:            cfg,
		userRepository: userRepository,
	}
}

// GetAll returns all users
func (service *DefaultUserService) GetAll() ([]*model.User, error) {
	users, err := service.userRepository.GetAll()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all users failed")
	}

	return users, errors.Wrap(err, "get all users failed")
}

// Get returns a user by it's username
func (service *DefaultUserService) Get(username string) (*model.User, error) {
	user, err := service.userRepository.Get(username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")
	}

	return user, errors.Wrap(err, "get user failed")
}

// GetByToken returns a user by token
func (service *DefaultUserService) GetByToken(rawToken string) (*model.User, error) {
	user, err := service.userRepository.GetByToken(rawToken)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user by token failed")
	}

	return user, errors.Wrap(err, "get user by token failed")
}

// Login checks whether the combination of username and password is valid
func (service *DefaultUserService) Login(username, password string) (*model.User, bool, error) {
	user, err := service.userRepository.Get(username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		return nil, false, err
	}

	return user, user != nil && user.Password == password, nil
}
