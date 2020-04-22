package db

import (
	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/dbx"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"
)

// IRepository is an interface of repository
type IRepository interface {}

// repository persists albums in database
type repository struct {
	db                db.IDB
	logger            log.ILogger
	defaultConditions map[string]interface{}
}

// Limit is default limit
const Limit = 100

// GetRepository return a repository
func GetRepository(dbase db.IDB, logger log.ILogger, entity string) (repo IRepository, err error) {
	r := &repository{
		db:     dbase,
		logger: logger,
	}

	switch entity {
	case "event":
		repo, err = NewEventRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}

