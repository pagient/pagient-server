package storage

import (
	"fmt"

	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
)

func New(cfg *config.Config) (model.Storage, error) {
	switch cfg.Database.Provider {
	case config.DatabaseProviderFile:
		return NewFileStorage(cfg.General.Root), nil
	default:
		return nil, fmt.Errorf("No suitable Database Provider found.")
	}
}
