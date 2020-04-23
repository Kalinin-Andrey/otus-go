package controller

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror/apperror"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror/grpcerror"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/calendarpb"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
)

// EventController is a handler for grpc requests
type EventController struct {
	Service event.IService
	Logger  log.ILogger
}

var _ calendarpb.CalendarServer = (*EventController)(nil)

// EventCreate creates an entity
func (c EventController) EventCreate(ctx context.Context, entityProto *calendarpb.Event) (*calendarpb.ResponseEvent, error) {
	var err error

	entity, err := EventProtoToEvent(*entityProto)
	if err != nil {
		return grpcerror.ErrBadRequestResponseEvent, apperror.ErrBadRequest
	}

	if err := entity.Validate(); err != nil {
		c.Logger.With(ctx).Info(err)
		return grpcerror.ErrBadRequestResponseEvent, apperror.ErrBadRequest
	}

	if err := c.Service.Create(ctx, entity); err != nil {
		c.Logger.With(ctx).Info(err)
		return grpcerror.ErrBadRequestResponseEvent, apperror.ErrBadRequest
	}

	entityProto, err = EventToEventProto(*entity)
	if err != nil {
		return grpcerror.ErrInternalResponseEvent, apperror.ErrInternal
	}

	return &calendarpb.ResponseEvent{
		Status:	grpcerror.StatusOK,
		Item:	entityProto,
	}, nil
}

// EventUpdate updates an entity
func (c EventController) EventUpdate(ctx context.Context, entityProto *calendarpb.Event) (*calendarpb.ResponseEvent, error) {
	var err error

	entity, err := EventProtoToEvent(*entityProto)
	if err != nil {
		return grpcerror.ErrBadRequestResponseEvent, apperror.ErrBadRequest
	}

	if err := entity.Validate(); err != nil {
		c.Logger.With(ctx).Info(err)
		return grpcerror.ErrBadRequestResponseEvent, apperror.ErrBadRequest
	}

	if err := c.Service.Update(ctx, entity); err != nil {
		c.Logger.With(ctx).Info(err)
		return grpcerror.ErrBadRequestResponseEvent, apperror.ErrBadRequest
	}

	entityProto, err = EventToEventProto(*entity)
	if err != nil {
		return grpcerror.ErrInternalResponseEvent, apperror.ErrInternal
	}

	return &calendarpb.ResponseEvent{
		Status:	grpcerror.StatusOK,
		Item:	entityProto,
	}, nil
}

// EventDelete deletes an entity
func (c EventController) EventDelete(ctx context.Context, idProto *calendarpb.EventID) (*calendarpb.Status, error) {
	id := uint(idProto.ID)

	if err := c.Service.Delete(ctx, id); err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx).Info(err)
			return grpcerror.StatusNotFound, apperror.ErrNotFound
		}
		c.Logger.With(ctx).Error(err)
		return grpcerror.StatusInternalServerError, apperror.ErrInternal
	}
	return grpcerror.StatusOK, nil
}

// EventListOnDay returns a list of events on a day
func (c EventController) EventListOnDay(ctx context.Context, ts *timestamp.Timestamp) (*calendarpb.ResponseEvents, error) {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		return grpcerror.ErrBadRequestResponseEvents, apperror.ErrBadRequest
	}
	return c.listFromTill(ctx, t, t.AddDate(0, 0, 1))
}

// EventListOnWeek returns a list of events on a week
func (c EventController) EventListOnWeek(ctx context.Context, ts *timestamp.Timestamp) (*calendarpb.ResponseEvents, error) {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		return grpcerror.ErrBadRequestResponseEvents, apperror.ErrBadRequest
	}
	return c.listFromTill(ctx, t, t.AddDate(0, 0, 7))
}

// EventListOnMonth returns a list of events on a month
func (c EventController) EventListOnMonth(ctx context.Context, ts *timestamp.Timestamp) (*calendarpb.ResponseEvents, error) {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		return grpcerror.ErrBadRequestResponseEvents, apperror.ErrBadRequest
	}
	return c.listFromTill(ctx, t, t.AddDate(0, 1, 0))
}

// listFromTill returns a list of events
func (c EventController) listFromTill(ctx context.Context, from time.Time, till time.Time) (*calendarpb.ResponseEvents, error) {

	condition := &event.QueryCondition{
		Where: &event.WhereCondition{
			Time: &event.WhereConditionTime{
				Between: &[2]time.Time{
					from,
					till,
				},
			},
		},
	}

	items, err := c.Service.List(ctx, condition)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx).Info(err)
			return grpcerror.ErrNotFoundResponseEvents, apperror.ErrNotFound
		}
		c.Logger.With(ctx).Error(err)
		return grpcerror.ErrInternalResponseEvents, apperror.ErrInternal
	}

	list, err := EventsListToEventProtosList(items)
	if err != nil {
		c.Logger.With(ctx).Error(err)
		return grpcerror.ErrInternalResponseEvents, apperror.ErrInternal
	}

	return &calendarpb.ResponseEvents{
		Status: &calendarpb.Status{
			OK:	true,
		},
		List:	list,
	}, nil
}



