package database

import (
	"database/sql"
	"fmt"
	"net/url"

	bridgeModel "github.com/pagient/pagient-server/internal/bridge/model"
	"github.com/pagient/pagient-server/internal/config"

	_ "github.com/denisenkom/go-mssqldb" // import mssql for database connection
	"github.com/pkg/errors"
)

// DB interface
type DB interface {
	GetRoomAssignments(string, ...uint) ([]*bridgeModel.RoomAssignment, error)
	Close() error
}

type db struct {
	*sql.DB
}

// Close closes the database
func (db *db) Close() error {
	return db.DB.Close()
}

// Open opens a sqlserver database connection
// uses global config for connection parameters
func Open() (DB, error) {
	if config.Bridge.DB.Driver != "sqlserver" {
		return nil, errors.New("only sqlserver is supported at the moment")
	}

	query := url.Values{}
	query.Add("database", config.Bridge.DB.Name)

	connURL := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(config.Bridge.DB.User, config.Bridge.DB.Password),
		Host:     fmt.Sprintf("%s:%d", config.Bridge.DB.Host, config.Bridge.DB.Port),
		RawQuery: query.Encode(),
	}

	dbConn, err := sql.Open(config.Bridge.DB.Driver, connURL.String())

	return &db{dbConn}, errors.Wrap(err, "could not connect to database server")
}
