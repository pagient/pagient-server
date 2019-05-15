package database

import (
	"log"
	"os"

	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // import sqlite for database connection
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// DB interface
type DB interface {
	Begin() (service.Tx, error)
	Close() error
}

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

// Open opens a sqlite3 database connection
// uses global config for connection parameters
func Open() (DB, error) {
	if config.DB.Driver != "sqlite3" {
		return nil, errors.New("only sqlite3 is supported at the moment")
	}

	dbConn, err := gorm.Open(config.DB.Driver, config.DB.Path)
	if err != nil {
		return nil, errors.New("establish database connection failed")
	}

	dbConn.LogMode(zerolog.GlobalLevel() <= zerolog.DebugLevel)
	dbConn.SetLogger(log.New(os.Stdout, "\r\n", 0))

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
