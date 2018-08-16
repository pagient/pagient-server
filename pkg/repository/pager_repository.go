package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

type pagerRepository struct {
	sqlRepository
}

// NewPagerRepository returns a new instance of a PagerRepository
func NewPagerRepository(db *gorm.DB) service.PagerRepository {
	return &pagerRepository{sqlRepository{db}}
}

// GetAll returns all available pagers
func (repo *pagerRepository) GetAll(sess service.DB) ([]*model.Pager, error) {
	session := sess.(*gorm.DB)

	var pagers []*model.Pager
	err := session.Find(&pagers).Error

	return pagers, errors.Wrap(err, "select all pagers failed")
}

// GetUnassigned returns all unassigned pagers
func (repo *pagerRepository) GetUnassigned(sess service.DB) ([]*model.Pager, error) {
	session := sess.(*gorm.DB)

	var pagers []*model.Pager
	err := session.
		Joins("LEFT JOIN patients ON patients.pager_id = pagers.id").
		Where("patients.id IS NULL").Find(&pagers).Error

	return pagers, errors.Wrap(err, "select unassigned pagers failed")
}

// Get returns a single pager by ID
func (repo *pagerRepository) Get(sess service.DB, id uint) (*model.Pager, error) {
	session := sess.(*gorm.DB)

	pager := &model.Pager{}
	err := session.First(pager, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return pager, errors.Wrap(err, "select pager by id failed")
}
