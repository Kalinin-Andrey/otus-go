package db

import (
	"fmt"
	"time"

	// pq is the driver for the postgres dialect
	_ "github.com/lib/pq"
	//_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"
)

// ConnectionTimeout is the default timeout for connection
const ConnectionTimeout = time.Duration(30 * time.Second)

// IDB is the interface for a DB connection
type IDB interface {
	DB() *sqlx.DB
}

// DB is the struct for a DB connection
type DB struct {
	db *sqlx.DB
}

// DB returns a db object
func (db *DB) DB() *sqlx.DB {
	return db.db
}

// Close connection
func (db *DB) Close() error {
	return db.db.Close()
}

var _ IDB = (*DB)(nil)

// New creates a new DB connection
func New(conf config.DB, logger log.ILogger) (*DB, error) {
	//db, err := gorm.Open(conf.Dialect, conf.DSN)
	db, err := ConnectLoop(conf.Dialect, conf.DSN, ConnectionTimeout)

	if err != nil {
		return nil, err
	}
	//db.SetLogger(logger)
	// Enable Logger, show detailed log
	//db.LogMode(true)

	dbobj := &DB{db: db}

	return dbobj, nil
}

// ConnectLoop is the func for connection in a loop with timeout
func ConnectLoop(dialect string, dsn string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := sqlx.Connect(dialect, dsn)
			if err == nil {
				return db, nil
			}
			//errors.Wrapf(err, "Can not connect to db %s by dsn: %q", dialect, dsn)
		}
	}
}

