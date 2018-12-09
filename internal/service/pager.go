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
