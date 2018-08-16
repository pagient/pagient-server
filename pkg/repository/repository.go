package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pkg/errors"
	"github.com/pagient/pagient-server/pkg/service"
)

type sqlRepository struct {
	db *gorm.DB
}

// BeginTx begins a transaction
func (repo *sqlRepository) BeginTx() service.DB {
	return repo.db.Begin()
}

//
func (repo *sqlRepository) RollbackTx(sess service.DB) service.DB {
	session := sess.(*gorm.DB)
	return session.Rollback()
}

func (repo *sqlRepository) CommitTx(sess service.DB) service.DB {
	session := sess.(*gorm.DB)
	return session.Commit()
}

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
