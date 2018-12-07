package service

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// GetAll returns all users
func (service *DefaultService) ListUsers() ([]*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	users, err := tx.GetUsers()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all users failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get all users failed")
	}

	tx.Commit()
	return users, nil
}

// Get returns a user by it's username
func (service *DefaultService) ShowUser(username string) (*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	user, err := tx.GetUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get user failed")
	}

	tx.Commit()
	return user, nil
}

// GetByToken returns a user by token
func (service *DefaultService) ShowUserByToken(rawToken string) (*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	user, err := tx.GetUserByToken(rawToken)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user by token failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get user by token failed")
	}

	tx.Commit()
	return user, nil
}

func (service *DefaultService) CreateUser(user *model.User) (*model.User, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	if err := service.validateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	user, err = tx.AddUser(user)
	if err != nil {
		log.Error().
			Err(err).
			Msg("add user failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "add user failed")
	}

	tx.Commit()
	return user, nil
}

// Login checks whether the combination of username and password is valid
func (service *DefaultService) Login(username, password string) (*model.User, bool, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, false, errors.Wrap(err, "create transaction failed")
	}

	user, err := tx.GetUser(username)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get user failed")

		tx.Rollback()
		return nil, false, errors.Wrap(err, "get user failed")
	}

	tx.Commit()
	return user, user != nil && comparePasswords(user.Password, password), nil
}

func (service *DefaultService) validateUser(tx Tx, user *model.User) error {
	var clients []*model.Client

	if user.ClientID != 0 {
		// load clients to validate if client sent with request is valid
		var err error
		clients, err = tx.GetClients()
		if err != nil {
			return errors.Wrap(err, "get all clients failed")
		}
	}

	// validate user
	if err := user.Validate(clients); err != nil {
		if model.IsValidationErr(err) {
			return &modelValidationErr{err.Error()}
		}

		return errors.Wrap(err, "validate user failed")
	}

	return nil
}

func hashPassword(plainPwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPwd), bcrypt.DefaultCost)
	return string(bytes), errors.Wrap(err, "generate bcrypt hash from password failed")
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
