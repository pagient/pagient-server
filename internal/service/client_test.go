package service

import (
	"testing"

	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultService_CreateClient(t *testing.T) {
	tests := map[string]struct {
		client        *model.Client
		validationErr error
		dbErr         error
	}{
		"successfully create client": {
			client: &model.Client{
				Name: "test",
			},
			validationErr: nil,
			dbErr:         nil,
		},
		"validate before create client": {
			client:        &model.Client{},
			validationErr: &modelValidationErr{},
			dbErr:         nil,
		},
		"rollback on db error": {
			client: &model.Client{
				Name: "test",
			},
			validationErr: nil,
			dbErr:         errors.New("test error"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		tx := &MockTx{}

		if test.validationErr == nil {
			txMethod := "Commit"
			if test.dbErr != nil {
				txMethod = "Rollback"
			}
			tx.On(txMethod).Return(nil).Once()

			tx.On("AddClient", mock.AnythingOfType("*model.Client")).Return(func(c *model.Client) error {
				c.ID = 1
				return test.dbErr
			}).Once()
		}

		db := &MockDB{}
		if test.validationErr == nil {
			db.On("Begin").Return(tx, nil).Once()
		}

		s := NewService(db, nil)
		err := s.CreateClient(test.client)

		id := uint(1)
		if test.validationErr != nil {
			id = uint(0)
		}
		assert.Equal(t, id, test.client.ID)
		if test.dbErr != nil {
			assert.EqualError(t, errors.Cause(err), test.dbErr.Error())
		}

		db.AssertExpectations(t)
		tx.AssertExpectations(t)
	}
}
