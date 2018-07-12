package repository

import (
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
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

func (repo *tokenFileRepository) Get(username string) ([]*model.Token, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	tokens := &[]*model.Token{}
	if err := repo.db.Read(tokenCollection, username, tokens); err != nil && !isNotFoundErr(err) {
		return nil, errors.Wrap(err, "read token failed")
	}

	return *tokens, nil
}

func (repo *tokenFileRepository) Add(token *model.Token) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	tokens := &[]*model.Token{}
	if err := repo.db.Read(tokenCollection, token.User, tokens); err != nil && !isNotFoundErr(err) {
		return errors.Wrap(err, "read token failed")
	}

	*tokens = append(*tokens, token)

	err := repo.db.Write(tokenCollection, token.User, tokens)
	return errors.Wrap(err, "write token failed")
}

func (repo *tokenFileRepository) Remove(token *model.Token) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	tokens := &[]*model.Token{}
	if err := repo.db.Read(tokenCollection, token.User, tokens); err != nil && !isNotFoundErr(err) {
		if isNotFoundErr(err) {
			return &entryNotExistErr{"token not found"}
		}
		return errors.Wrap(err, "read token failed")
	}

	for i, tok := range *tokens {
		if tok.Token == token.Token {
			*tokens = append((*tokens)[:i], (*tokens)[i+1:]...)
			break
		}
	}

	var err error
	if len(*tokens) == 0 {
		err = repo.db.Delete(tokenCollection, token.User)
	} else {
		err = repo.db.Write(tokenCollection, token.User, tokens)
	}

	return errors.Wrap(err, "delete token failed")
}
