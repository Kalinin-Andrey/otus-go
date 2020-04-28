package event

import (
	"context"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/notification"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror/apperror"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"
)

// MaxLIstLimit const
const MaxLIstLimit = 1000

// IService encapsulates usecase logic for event.
type IService interface {
	// NewEntity returns new empty entity
	NewEntity() *Event
	// Get returns an entity with given ID
	Get(ctx context.Context, id uint) (*Event, error)
	//First(ctx context.Context, entity *Event) (*Event, error)
	// Query returns a list with pagination
	Query(ctx context.Context, offset, limit uint) ([]Event, error)
	// List entity
	List(ctx context.Context, condition *QueryCondition) ([]Event, error)
	//Count(ctx context.Context) (uint, error)
	// Create entity
	Create(ctx context.Context, entity *Event) error
	// Update entity
	Update(ctx context.Context, entity *Event) error
	// Delete entity
	Delete(ctx context.Context, id uint) (error)
	// NotificationsList returns a list of notification.
	NotificationsList(ctx context.Context) ([]*notification.Notification, error)
}

type service struct {
	//Domain     Domain
	repo       			IRepository
	logger     			log.ILogger
}

// NewService creates a new service.
func NewService(repo IRepository, logger log.ILogger) IService {
	s := &service{repo, logger}
	//repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s service) defaultConditions() map[string]interface{} {
	return map[string]interface{}{
	}
}

// NewEntity returns a new empty entity
func (s service) NewEntity() *Event {
	return &Event{}
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*Event, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		if err == apperror.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "Can not get a event by id: %v", id)
	}
	return entity, nil
}

// First returns an entity representing one new record
/*func (s service) First(ctx context.Context, entity *Event) (*Event, error) {
	return s.repo.First(ctx, entity)
}*/

/*
// Count returns the number of items.
func (s service) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit uint) ([]Event, error) {
	items, err := s.repo.Query(ctx, nil, offset, limit)
	if err != nil {
		if err == apperror.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "Can not find a list of events by ctx")
	}
	return items, nil
}


// List returns the items list.
func (s service) List(ctx context.Context, query *QueryCondition) ([]Event, error) {
	items, err := s.repo.Query(ctx, query, 0, MaxLIstLimit)
	if err != nil {
		if err == apperror.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "Can not find a list of events by ctx")
	}
	return items, nil
}

// Create entity
func (s service) Create(ctx context.Context, entity *Event) error {
	return s.repo.Create(ctx, entity)
}

// Update entity
func (s service) Update(ctx context.Context, entity *Event) error {
	return s.repo.Update(ctx, entity)
}

// Delete entity
func (s service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// NotificationsList returns a list of notification.
func (s service) NotificationsList(ctx context.Context) ([]*notification.Notification, error) {
	items, err := s.repo.ListForNotifications(ctx, 0, MaxLIstLimit)
	if err != nil {
		if err == apperror.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "Can not find a list of events by ctx")
	}
	return s.NotificationsForEvents(items), nil
}

// NotificationsForEvents returns a list of notifications for a given list of events
func (s service) NotificationsForEvents(eList []Event) []*notification.Notification {
	nList := make([]*notification.Notification, 0, len(eList))

	for _, e := range eList {
		nList = append(nList, e.Notification())
	}

	return nList
}

