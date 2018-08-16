package service

import (
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// PagerService interface
type PagerService interface {
	GetAll() ([]*model.Pager, error)
	Get(uint) (*model.Pager, error)
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
	session := service.pagerRepository.BeginTx()
	pagers, err := service.pagerRepository.GetAll(session)
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all pagers failed")

		service.pagerRepository.RollbackTx(session)
		return nil, errors.Wrap(err, "get all pagers failed")
	}

	service.pagerRepository.CommitTx(session)
	return pagers, nil
}

// Get returns a pager by it's id
func (service *DefaultPagerService) Get(id uint) (*model.Pager, error) {
	session := service.pagerRepository.BeginTx()
	pager, err := service.pagerRepository.Get(session, id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("pager ID", id).
			Msg("get pager failed")

		service.pagerRepository.RollbackTx(session)
		return nil, errors.Wrapf(err, "get pager failed", id)
	}

	service.pagerRepository.CommitTx(session)
	return pager, nil
}
