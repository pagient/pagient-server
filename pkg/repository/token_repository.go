package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/pagient/pagient-server/pkg/service"
)

type tokenRepository struct {
	sqlRepository
}

// NewTokenRepository returns a new instance of a TokenRepository
func NewTokenRepository(db *gorm.DB) service.TokenRepository {
	return &tokenRepository{sqlRepository{db}}
}

func (repo *tokenRepository) Get(sess service.DB, rawToken string) (*model.Token, error) {
	session := sess.(*gorm.DB)

	token := &model.Token{}
	err := session.Where(&model.Token{
		Raw: rawToken,
	}).First(token).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return token, errors.Wrap(err, "select token failed")
}

func (repo *tokenRepository) GetByUser(sess service.DB, username string) ([]*model.Token, error) {
	session := sess.(*gorm.DB)

	var tokens []*model.Token
	err := session.
		Joins("JOIN users ON users.id = tokens.user_id").
		Where("users.username = ?", username).Find(&tokens).Error

	return tokens, errors.Wrap(err, "select tokens by user failed")
}

func (repo *tokenRepository) Add(sess service.DB, token *model.Token) (*model.Token, error) {
	session := sess.(*gorm.DB)

	err := session.Create(token).Error

	return token, errors.Wrap(err, "create token failed")
}

func (repo *tokenRepository) Remove(sess service.DB, token *model.Token) (*model.Token, error) {
	session := sess.(*gorm.DB)

	err := session.Delete(token).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, &entryNotExistErr{"token not found"}
	}

	return token, errors.Wrap(err, "delete token failed")
}
