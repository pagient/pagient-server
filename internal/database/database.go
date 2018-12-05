package database

import (
	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // import sqlite for database connection
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type db struct {
	*gorm.DB
}

// Begin starts an returns a new transaction.
func (db *db) Begin() (service.Tx, error) {
	t := db.DB.Begin()
	return &tx{t}, errors.Wrap(t.Error, "begin new transaction failed")
}

// Close closes the database
func (db *db) Close() error {
	return db.DB.Close()
}

type tx struct {
	*gorm.DB
}

func (t *tx) Commit() error {
	return t.DB.Commit().Error
}

func (t *tx) Rollback() error {
	return t.DB.Rollback().Error
}

func Open() (*db, error) {
	if config.General.DB.Driver != "sqlite3" {
		return nil, errors.New("only sqlite3 is supported at the moment")
	}

	dbConn, err := gorm.Open(config.General.DB.Driver, config.DB.Path)
	if err != nil {
		return nil, errors.New("establish database connection failed")
	}

	dbConn.LogMode(zerolog.GlobalLevel() <= zerolog.DebugLevel)
	dbConn.SetLogger(&log.Logger)

	// Create database tables etc.
	if err := createTables(dbConn); err != nil {
		return nil, errors.New("create database tables failed")
	}

	return &db{dbConn}, nil
}

// creates necessary database tables
func createTables(db *gorm.DB) error {
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
