package database

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetToken returns the token by it's raw string form
func (t *tx) GetToken(rawToken string) (*model.Token, error) {
	token := &model.Token{}
	err := t.Where(&model.Token{
		Raw: rawToken,
	}).First(token).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return token, errors.Wrap(err, "select token failed")
}

// GetTokensByUser returns all tokens from a user
func (t *tx) GetTokensByUser(username string) ([]*model.Token, error) {
	var tokens []*model.Token
	err := t.Joins("JOIN users ON users.id = tokens.user_id").
		Where("users.username = ?", username).Find(&tokens).Error

	return tokens, errors.Wrap(err, "select tokens by user failed")
}

// AddToken adds a token
func (t *tx) AddToken(token *model.Token) (*model.Token, error) {
	err := t.Create(token).Error

	return token, errors.Wrap(err, "create token failed")
}

// RemoveToken removes a token
func (t *tx) RemoveToken(token *model.Token) error {
	err := t.Delete(token).Error
	if gorm.IsRecordNotFoundError(err) {
		return &entryNotExistErr{"token not found"}
	}

	return errors.Wrap(err, "delete token failed")
}
