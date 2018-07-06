package repository

import (
	"strconv"
	"strings"
	"sync"

	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/pkg/errors"
)

var (
	pagerRepositoryOnce     sync.Once
	pagerRepositoryInstance service.PagerRepository
)

// GetPagerRepositoryInstance creates and returns a new PagerCfgRepository
func GetPagerRepositoryInstance(cfg *config.Config) (service.PagerRepository, error) {
	pagerRepositoryOnce.Do(func() {
		pagerRepositoryInstance = &pagerCfgRepository{cfg}
	})

	return pagerRepositoryInstance, nil
}

type pagerCfgRepository struct {
	cfg *config.Config
}

// GetAll returns all available pagers
func (repo *pagerCfgRepository) GetAll() ([]*model.Pager, error) {
	var pagers []*model.Pager
	for _, pagerInfo := range repo.cfg.General.Pagers {
		pair := strings.SplitN(pagerInfo, ":", 2)

		id, err := strconv.Atoi(pair[0])
		if err != nil {
			return nil, errors.Wrap(err, "integer string conversion failed")
		}

		pagers = append(pagers, &model.Pager{ID: id, Name: pair[1]})
	}

	return pagers, nil
}

// Get returns a single pager by ID
func (repo *pagerCfgRepository) Get(id int) (*model.Pager, error) {
	pagers, err := repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "get pagers failed")
	}

	for _, pager := range pagers {
		if pager.ID == id {
			return pager, nil
		}
	}

	return nil, nil
}
