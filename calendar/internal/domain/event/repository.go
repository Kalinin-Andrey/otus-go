package event

import(

)

// IEventRepository is a interface for events repositories
type IEventRepository interface {
	Create(event *Event) (uint, error)
	Read(eventID uint) (*Event, error)
	ReadAll() (map[uint]*Event, error)
	Update(event *Event) error
	Delete(eventID uint) error
	IsTimeNotBusy(event *Event) bool
}


