package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/pagient/pagient-server/internal/bridge"
	"github.com/pagient/pagient-server/internal/config"

	_ "github.com/denisenkom/go-mssqldb" // import mssql for database connection
	"github.com/pkg/errors"
)

type db struct {
	*sql.DB
}

// Begin starts an returns a new transaction.
func (db *db) Begin() (bridge.Tx, error) {
	t, err := db.DB.Begin()
	return &tx{t}, errors.Wrap(err, "begin new transaction failed")
}

// Close closes the database
func (db *db) Close() error {
	return db.DB.Close()
}

type tx struct {
	*sql.Tx
}

func (t *tx) Commit() error {
	return t.Tx.Commit()
}

func (t *tx) Rollback() error {
	return t.Tx.Rollback()
}

// OpenSQL opens a mssql database connection by given config
func Open() (bridge.DB, error) {
	if config.Bridge.DB.Driver != "sqlserver" {
		return nil, errors.New("only sqlserver is supported at the moment")
	}

	query := url.Values{}
	query.Add("database", config.Bridge.DB.Name)
	query.Add("encrypt", "disable")

	connUrl := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(config.Bridge.DB.User, config.Bridge.DB.Password),
		Host:     fmt.Sprintf("%s:%d", config.Bridge.DB.Host, config.Bridge.DB.Port),
		RawQuery: query.Encode(),
	}

	dbConn, err := sql.Open(config.Bridge.DB.Driver, connUrl.String())

	return &db{dbConn}, errors.Wrap(err, "could not connect to database server")
}
