package mock

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
)

// EventRepository is a mock for EventRepository
type EventRepository struct {
	Response struct {
		Get		struct {
			Entity	*event.Event
			Err		error
		}
		First	struct {
			Entity	*event.Event
			Err		error
		}
		Query	struct {
			List	[]event.Event
			Err		error
		}
		Create	struct {
			Entity	*event.Event
			Err		error
		}
		Update	struct {
			Entity	*event.Event
			Err		error
		}
		Delete	struct {
			Err		error
		}
	}
}

var _ event.IRepository = (*EventRepository)(nil)

// SetDefaultConditions mock
func (r EventRepository) SetDefaultConditions(conditions map[string]interface{}) {}

// Get mock
func (r EventRepository) Get(ctx context.Context, id uint) (*event.Event, error) {
	return r.Response.Get.Entity, r.Response.Get.Err
}

// First mock
func (r EventRepository) First(ctx context.Context, user *event.Event) (*event.Event, error) {
	return r.Response.First.Entity, r.Response.First.Err
}

// Query mock
func (r EventRepository) Query(ctx context.Context, offset, limit uint) ([]event.Event, error) {
	return r.Response.Query.List, r.Response.Query.Err
}

// Create mock
func (r EventRepository) Create(ctx context.Context, entity *event.Event) error {
	if r.Response.Create.Entity != nil {
		copier.Copy(&entity, &r.Response.Create.Entity)
	}
	return r.Response.Create.Err
}

// Update mock
func (r EventRepository) Update(ctx context.Context, entity *event.Event) error {
	if r.Response.Create.Entity != nil {
		copier.Copy(&entity, &r.Response.Create.Entity)
	}
	return r.Response.Update.Err
}

// Delete mock
func (r EventRepository) Delete(ctx context.Context, id uint) error {
	return r.Response.Delete.Err
}
