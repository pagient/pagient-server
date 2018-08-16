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
	session := service.userRepository.BeginTx()
	users, err := service.userRepository.GetAll(session)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all users failed")

		service.userRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get all users failed")
	}

	service.userRepository.CommitTx(session)
	return users, nil
}

// Get returns a user by it's username
func (service *DefaultUserService) Get(username string) (*model.User, error) {
	session := service.userRepository.BeginTx()
	user, err := service.userRepository.Get(session, username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		service.userRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get user failed")
	}

	service.userRepository.CommitTx(session)
	return user, nil
}

// GetByToken returns a user by token
func (service *DefaultUserService) GetByToken(rawToken string) (*model.User, error) {
	session := service.userRepository.BeginTx()
	user, err := service.userRepository.GetByToken(session, rawToken)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user by token failed")

		service.userRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get user by token failed")
	}

	service.userRepository.CommitTx(session)
	return user, nil
}

// Login checks whether the combination of username and password is valid
func (service *DefaultUserService) Login(username, password string) (*model.User, bool, error) {
	session := service.userRepository.BeginTx()
	user, err := service.userRepository.Get(session, username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		service.userRepository.RollbackTx(session)
		return nil, false, errors.Wrap(err, "get user failed")
	}

	service.userRepository.CommitTx(session)
	return user, user != nil && user.Password == password, nil
}
