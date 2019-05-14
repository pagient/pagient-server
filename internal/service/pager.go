package service

import (
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ListPagers returns all pagers
func (service *defaultService) ListPagers() ([]*model.Pager, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	pagers, err := tx.GetPagers()
	if err != nil {
		log.Error().
			Err(err).
			Msg("get all pagers failed")

		tx.Rollback()
		return nil, errors.Wrap(err, "get all pagers failed")
	}

	tx.Commit()
	return pagers, nil
}

// ShowPager returns a pager by it's id
func (service *defaultService) ShowPager(id uint) (*model.Pager, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	pager, err := tx.GetPager(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("pager ID", id).
			Msg("get pager failed")

		tx.Rollback()
		return nil, errors.Wrapf(err, "get pager %d failed", id)
	}

	tx.Commit()
	return pager, nil
}

// CreatePager creates a new pager
func (service *defaultService) CreatePager(pager *model.Pager) error {
	if err := service.validatePager(pager); err != nil {
		return errors.WithStack(err)
	}

	tx, err := service.db.Begin()
	if err != nil {
		return errors.Wrap(err, "create transaction failed")
	}

	err = tx.AddPager(pager)
	if err != nil {
		log.Error().
			Err(err).
			Msg("add pager failed")

		tx.Rollback()
		return errors.Wrap(err, "add pager failed")
	}

	tx.Commit()
	return nil
}

func (service *defaultService) validatePager(pager *model.Pager) error {
	if err := pager.Validate(); err != nil {
		if model.IsValidationErr(err) {
			return &modelValidationErr{err.Error()}
		}

		return errors.Wrap(err, "validate pager failed")
	}

	return nil
}
