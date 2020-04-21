package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"time"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

)

// IDB is the interface for a DB connection
type IDB interface {
	DB() *gorm.DB
}

// DB is the struct for a DB connection
type DB struct {
	db *gorm.DB
}

// DB returns a db object
func (db *DB) DB() *gorm.DB {
	return db.db
}

var _ IDB = (*DB)(nil)

// New creates a new DB connection
func New(conf config.DB, logger log.ILogger) (*DB, error) {
	//db, err := gorm.Open(conf.Dialect, conf.DSN)
	db, err := ConnectLoop(conf.Dialect, conf.DSN, time.Duration(10 * time.Second))

	if err != nil {
		return nil, err
	}
	db.SetLogger(logger)
	// Enable Logger, show detailed log
	db.LogMode(true)
	// Enable auto preload embeded entities
	db = 	db.Set("gorm:auto_preload", true)

	dbobj := &DB{db: db}
	dbobj.AutoMigrateAll()
	
	return dbobj, nil
}

// ConnectLoop is the func for connection in a loop with timeout
func ConnectLoop(dialect string, dsn string, timeout time.Duration) (*gorm.DB, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := gorm.Open(dialect, dsn)
			if err == nil {
				return db, nil
			}
			fmt.Println(errors.Wrapf(err, "failed to connect to db %s", dsn))
		}
	}
}

