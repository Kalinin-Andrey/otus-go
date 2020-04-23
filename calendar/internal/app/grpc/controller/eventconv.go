package controller

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/calendarpb"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/golang/protobuf/ptypes"
)

// EventToEventProto convs Event To EventProto
func EventToEventProto(entity event.Event) (*calendarpb.Event, error) {
	time, err := ptypes.TimestampProto(entity.Time)
	if err != nil {
		return nil, err
	}
	duration := ptypes.DurationProto(entity.Duration)
	entityProto := &calendarpb.Event{
		Id:       uint64(entity.ID),
		UserId:   uint64(entity.UserID),
		Title:    entity.Title,
		Time:     time,
		Duration: duration,
	}
	if entity.Description != nil {
		entityProto.Description = *entity.Description
	}
	if entity.NoticePeriod != nil {
		noticePeriod := ptypes.DurationProto(*entity.NoticePeriod)
		entityProto.NoticePeriod = noticePeriod
	}
	return entityProto, nil
}

// EventProtoToEvent convs EventProto To Event
func EventProtoToEvent(entityProto calendarpb.Event) (*event.Event, error) {
	var err error
	entity	:= &event.Event{}
	entity.ID		= uint(entityProto.Id)
	entity.UserID	= uint(entityProto.UserId)
	entity.Title	= entityProto.Title
	if entityProto.Description != "" {
		entity.Description = &entityProto.Description
	}
	if entityProto.Time != nil {
		entity.Time, err = ptypes.Timestamp(entityProto.Time)
		if err != nil {
			return nil, err
		}
	}
	if entityProto.Duration != nil {
		entity.Duration, err = ptypes.Duration(entityProto.Duration)
		if err != nil {
			return nil, err
		}
	}
	if entityProto.NoticePeriod != nil {
		noticePeriod, err := ptypes.Duration(entityProto.NoticePeriod)
		if err != nil {
			return nil, err
		}
		entity.NoticePeriod = &noticePeriod
	}
	return entity, nil
}

// EventsListToEventProtosList convs Events List To EventProtos List
func EventsListToEventProtosList(items []event.Event) ([]*calendarpb.Event, error) {
	list := make([]*calendarpb.Event, 0, len(items))
	for _, item := range items {
		e, err := EventToEventProto(item)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}
