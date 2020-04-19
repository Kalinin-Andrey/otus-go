package db

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/user"
)

// AutoMigrateAll is the func for auto migration for all entities
func (db *DB) AutoMigrateAll() {
	db.DB().AutoMigrate(
		&user.User{},
		&event.Event{},
		)
}
