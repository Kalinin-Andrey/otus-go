package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/db"
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
	case "user":
		repo, err = NewUserRepository(r)
	case "event":
		repo, err = NewEventRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}

// SetDefaultConditions func
func  (r *repository) SetDefaultConditions(conditions map[string]interface{}) {
	r.defaultConditions = conditions

	if _, ok := r.defaultConditions["Limit"]; !ok {
		r.defaultConditions["Limit"] = Limit
	}
}

func (r repository) dbWithDefaults() *gorm.DB {
	db := r.db.DB()

	if where, ok := r.defaultConditions["Where"]; ok {
		m := r.keysToSnakeCase(where.(map[string]interface{}))
		db = db.Where(m)
	}

	if order, ok := r.defaultConditions["SortOrder"]; ok {
		m := r.keysToSnakeCase(order.(map[string]interface{}))
		db = db.Order(m)
	}

	if limit, ok := r.defaultConditions["Limit"]; ok {
		db = db.Limit(limit)
	}

	return db
}


func (r repository) dbWithContext(ctx context.Context, db *gorm.DB) *gorm.DB {

	if where := ctx.Value("Where"); where != nil {
		m := r.keysToSnakeCase(where.(map[string]interface{}))
		db = db.Where(m)
	}

	if order := ctx.Value("SortOrder"); order != nil {
		m := r.keysToSnakeCase(order.(map[string]interface{}))
		db = db.Order(m)
	}

	if limit := ctx.Value("Limit"); limit != nil {
		db = db.Limit(limit)
	}

	return db
}

func (r repository) keysToSnakeCase(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for key, val := range in {
		out[strcase.ToSnake(key)] = val
	}
	return out
}