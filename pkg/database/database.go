package database

import (
	"fmt"

	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/database/filedb"
	"github.com/pagient/pagient-api/pkg/model"
)

// New creates a new database connection based on configured provider
func New(cfg *config.Config) (model.Database, error) {
	switch cfg.Database.Provider {
	case config.DatabaseProviderFile:
		return filedb.New(cfg.General.Root)
	default:
		return nil, fmt.Errorf("no suitable database provider found")
	}
}
