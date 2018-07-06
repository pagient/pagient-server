package service

import (
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// PagerService interface
type PagerService interface {
	GetAll() ([]*model.Pager, error)
	Get(int) (*model.Pager, error)
}

// DefaultPagerService struct
type DefaultPagerService struct {
	pagerRepository PagerRepository
}

// NewPagerService initializes a PagerService
func NewPagerService(repository PagerRepository) PagerService {
	return &DefaultPagerService{
		pagerRepository: repository,
	}
}

// GetAll returns all pagers
func (service *DefaultPagerService) GetAll() ([]*model.Pager, error) {
	pagers, err := service.pagerRepository.GetAll()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all pagers failed")
	}

	return pagers, errors.Wrap(err, "get all pagers failed")
}

// Get returns a pager by it's id
func (service *DefaultPagerService) Get(id int) (*model.Pager, error) {
	pager, err := service.pagerRepository.Get(id)
	if err != nil {
		log.Error().
			Err(err).
			Int("pager ID", id).
			Msg("get pager failed")
	}

	return pager, errors.Wrapf(err, "get pager failed", id)
}
