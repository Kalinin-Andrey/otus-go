package db

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	//_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror/apperror"
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
	//var id int64
	entity := &event.Event{}

	//row := r.db.DB().QueryRowContext(ctx, "SELECT * FROM event WHERE id = $1", id)
	//err := row.Scan(&entity.ID, &entity.UserID, &entity.Title, &entity.Description, &entity.Time, &entity.Duration, &entity.NoticePeriod)
	err := r.db.DB().GetContext(ctx, entity, "SELECT * FROM event WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves records with the specified offset and limit from the database.
func (r EventRepository) Query(ctx context.Context, query *event.QueryCondition, offset, limit uint) ([]event.Event, error) {
	var items []event.Event

	if query == nil || query.Where == nil || query.Where.Time == nil || query.Where.Time.Between == nil {
		return nil, errors.Errorf("invalid query param: %v", query)
	}

	err := r.db.DB().SelectContext(ctx, &items, "SELECT * FROM event WHERE time BETWEEN $1 AND $2", query.Where.Time.Between[0], query.Where.Time.Between[1])
	if err != nil {
		if err == sql.ErrNoRows {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

// Create saves a new entity in the database.
func (r EventRepository) Create(ctx context.Context, entity *event.Event) error {
	var lastInsertID uint
	err := r.db.DB().QueryRowContext(ctx, "INSERT INTO event (user_id, title, description, \"time\", duration, notice_period) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", entity.UserID, entity.Title, entity.Description, entity.Time, entity.Duration, entity.NoticePeriod).Scan(&lastInsertID)
	if err != nil {
		return errors.Wrapf(err, "EventRepository: error inserting entity %v", entity)
	}
	/*id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrapf(err, "EventRepository: error on LastInsertId() result: %v", res)
	}*/
	entity.ID = lastInsertID
	return nil
}

// Update recoprd of entity in db
func (r EventRepository) Update(ctx context.Context, entity *event.Event) error {
	_, err := r.db.DB().ExecContext(ctx, "UPDATE event SET user_id = $1, title = $2, description = $3, \"time\" = $4, duration = $5, notice_period = $6", entity.UserID, entity.Title, entity.Description, entity.Time, entity.Duration, entity.NoticePeriod)
	if err != nil {
		return errors.Wrapf(err, "EventRepository: error updating entity %v", entity)
	}
	return nil
}

// Delete deletes a record with the specified ID from the database.
func (r EventRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.DB().ExecContext(ctx, "DELETE FROM event WHERE id = $1", id)
	if err != nil {
		return errors.Wrapf(err, "EventRepository: error deleting record id = %v", id)
	}
	return nil
}


// ListForNotifications returns list of events for notification
func (r EventRepository) ListForNotifications(ctx context.Context, offset, limit uint) ([]event.Event, error) {
	var items []event.Event

	// Now() >= e.Time - e.NoticePeriod
	err := r.db.DB().SelectContext(ctx, &items, "SELECT * FROM event WHERE \"time\" -  make_interval(0, 0, 0, 0, 0, 0, notice_period / 1000000000) <= now()")
	if err != nil {
		if err == sql.ErrNoRows {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

