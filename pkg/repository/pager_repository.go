package repository

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
)

var (
	pagerRepositoryOnce     sync.Once
	pagerRepositoryInstance service.PagerRepository
)

// GetPagerRepositoryInstance creates and returns a new PagerCfgRepository
func GetPagerRepositoryInstance(db *gorm.DB) (service.PagerRepository, error) {
	pagerRepositoryOnce.Do(func() {
		pagerRepositoryInstance = &pagerRepository{db}
	})

	return pagerRepositoryInstance, nil
}

type pagerRepository struct {
	db *gorm.DB
}

// GetAll returns all available pagers
func (repo *pagerRepository) GetAll() ([]*model.Pager, error) {
	var pagers []*model.Pager
	err := repo.db.Find(&pagers).Error

	return pagers, errors.Wrap(err, "select all pagers failed")
}

// GetUnassigned returns all unassigned pagers
func (repo *pagerRepository) GetUnassigned() ([]*model.Pager, error) {
	var pagers []*model.Pager
	err := repo.db.
		Joins("LEFT JOIN patients ON patients.pager_id = pagers.id").
		Where("patients.id IS NULL").Find(&pagers).Error

	return pagers, errors.Wrap(err, "select unassigned pagers failed")
}

// Get returns a single pager by ID
func (repo *pagerRepository) Get(id uint) (*model.Pager, error) {
	pager := &model.Pager{}
	err := repo.db.First(pager, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return pager, errors.Wrap(err, "select pager by id failed")
}
