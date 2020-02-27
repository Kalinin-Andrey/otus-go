package db

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"
	"github.com/go-ozzo/ozzo-dbx"
	//"database/sql"
	_ "github.com/lib/pq"
)

// IDB is the interface for a DB connection
type IDB interface {

}

// DB is the struct for a DB connection
type DB struct {
	db	*dbx.DB
}

// New creates a new DB connection
func New(conf config.DB, logger log.ILogger) (*DB, error) {
	db, err := dbx.MustOpen("postgres", conf.DSN)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()
	return NewWithDbx(db), nil
}

// NewWithDbx creates a new DB connection with dbx.DB
func NewWithDbx(db *dbx.DB) *DB {
	return &DB{db}
}




