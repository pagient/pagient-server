package repository

import (
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/pkg/errors"
)

const (
	tokenCollection = "token"
)

var (
	tokenRepositoryOnce     sync.Once
	tokenRepositoryInstance service.TokenRepository
)

// GetTokenRepositoryInstance creates and returns a new TokenFileRepository
func GetTokenRepositoryInstance(cfg *config.Config) (service.TokenRepository, error) {
	var err error

	tokenRepositoryOnce.Do(func() {
		// Set up scribble json file store
		var db fileDriver
		db, err = scribble.New(cfg.General.Root, nil)

		tokenRepositoryInstance = &tokenFileRepository{
			lock: &sync.Mutex{},
			db:   db,
		}
	})

	if err != nil {
		return nil, errors.Wrap(err, "init scribble store failed")
	}

	return tokenRepositoryInstance, nil
}

type tokenFileRepository struct {
	lock *sync.Mutex
	db   fileDriver
}

func (repo *tokenFileRepository) Get(username string) (*model.Token, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	token := &model.Token{}
	if err := repo.db.Read(tokenCollection, username, token); err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "read token failed")
	}

	return token, nil
}

func (repo *tokenFileRepository) Add(username string, token *model.Token) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	tok := &model.Token{}
	if err := repo.db.Read(tokenCollection, username, tok); err != nil && !isNotFoundErr(err) {
		return errors.Wrap(err, "read token failed")
	}
	if tok.Token != "" {
		return &entryExistErr{"token already exists"}
	}

	err := repo.db.Write(tokenCollection, username, token)
	return errors.Wrap(err, "write token failed")
}

func (repo *tokenFileRepository) Remove(username string) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	err := repo.db.Delete(tokenCollection, username)
	if err != nil {
		if isNotFoundErr(err) {
			return &entryNotExistErr{"token not found"}
		}
		return errors.Wrap(err, "delete token failed")
	}

	return nil
}


