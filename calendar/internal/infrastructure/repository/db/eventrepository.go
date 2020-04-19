package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror"
)

// EventRepository is a repository for the event entity
type EventRepository struct {
	repository
}

var _ event.IRepository = (*EventRepository)(nil)

// NewEventRepository creates a new Repository
func NewEventRepository(repository *repository) (*EventRepository, error) {
	return &EventRepository{repository: *repository}, nil
}


// Get reads entities with the specified ID from the database.
func (r EventRepository) Get(ctx context.Context, id uint) (*event.Event, error) {
	entity := &event.Event{}

	err := r.dbWithDefaults().First(&entity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// First returns an entity representing a one first record
func (r EventRepository) First(ctx context.Context, entity *event.Event) (*event.Event, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves records with the specified offset and limit from the database.
func (r EventRepository) Query(ctx context.Context, offset, limit uint) ([]event.Event, error) {
	var items []event.Event

	err := r.dbWithContext(ctx, r.dbWithDefaults()).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

// Create saves a new entity in the database.
func (r EventRepository) Create(ctx context.Context, entity *event.Event) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update recoprd of entity in db
func (r EventRepository) Update(ctx context.Context, entity *event.Event) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}
	return r.db.DB().Save(entity).Error
}

// Delete deletes a record with the specified ID from the database.
func (r EventRepository) Delete(ctx context.Context, id uint) error {
	entity, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.DB().Delete(entity).Error
}

