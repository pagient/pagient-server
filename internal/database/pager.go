package database

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetAll returns all available pagers
func (t *tx) GetPagers() ([]*model.Pager, error) {
	var pagers []*model.Pager
	err := t.Find(&pagers).Error

	return pagers, errors.Wrap(err, "select all pagers failed")
}

// GetUnassigned returns all unassigned pagers
func (t *tx) GetUnassignedPagers() ([]*model.Pager, error) {
	var pagers []*model.Pager
	err := t.Joins("LEFT JOIN patients ON patients.pager_id = pagers.id").
		Where("patients.id IS NULL").Find(&pagers).Error

	return pagers, errors.Wrap(err, "select unassigned pagers failed")
}

// Get returns a single pager by ID
func (t *tx) GetPager(id uint) (*model.Pager, error) {
	pager := &model.Pager{}
	err := t.First(pager, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return pager, errors.Wrap(err, "select pager by id failed")
}
