package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
)

// InitDatabase creates necessary database tables
func InitDatabase(db *gorm.DB) error {
	tables := []interface{}{
		&model.Client{},
		&model.Pager{},
		&model.Patient{},
		&model.Token{},
		&model.User{},
	}

	for _, table := range tables {
		if !db.HasTable(table) {
			if err := db.CreateTable(table).Error; err != nil {
				return errors.Wrap(err, "create table failed")
			}
		}
	}

	return nil
}
