package event

import(

)

type IEventRepository interface {
	Create(event *Event) error
	Read(eventId uint) (*Event, error)
	ReadAll() (map[uint]*Event, error)
	Update(event *Event) error
	Delete(eventId uint) error
}


