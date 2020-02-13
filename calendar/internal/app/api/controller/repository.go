package event

import(

)

type IEventRepository interface {
	Create(event Event) error
	ReadOne(eventId uint) (Event, error)
	ReadAll() ([]Event, error)
	Update(event Event) error
	Delete(eventId uint) error
}


