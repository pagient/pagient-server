package repository

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

var (
	tokenRepositoryOnce     sync.Once
	tokenRepositoryInstance service.TokenRepository
)

// GetTokenRepositoryInstance creates and returns a new TokenFileRepository
func GetTokenRepositoryInstance(db *gorm.DB) (service.TokenRepository, error) {
	tokenRepositoryOnce.Do(func() {
		tokenRepositoryInstance = &tokenRepository{db}
	})

	return tokenRepositoryInstance, nil
}

type tokenRepository struct {
	db *gorm.DB
}

func (repo *tokenRepository) Get(rawToken string) (*model.Token, error) {
	token := &model.Token{}
	err := repo.db.Where(&model.Token{
		Raw: rawToken,
	}).First(token).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return token, errors.Wrap(err, "select token failed")
}

func (repo *tokenRepository) GetByUser(username string) ([]*model.Token, error) {
	var tokens []*model.Token
	err := repo.db.
		Joins("JOIN users ON users.id = tokens.user_id").
		Where("users.username = ?", username).Find(&tokens).Error

	return tokens, errors.Wrap(err, "select tokens by user failed")
}

func (repo *tokenRepository) Add(token *model.Token) (*model.Token, error) {
	err := repo.db.Create(token).Error

	return token, errors.Wrap(err, "create token failed")
}

func (repo *tokenRepository) Remove(token *model.Token) (*model.Token, error) {
	err := repo.db.Delete(token).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, &entryNotExistErr{"token not found"}
	}

	return token, errors.Wrap(err, "delete token failed")
}
