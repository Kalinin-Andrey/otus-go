package in_memory

import(
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/calendar_error"
)

const startCapacity = 100

type EventRepository struct {
	eventsList map[uint]*event.Event
}

var lastId uint

func nextId() uint {
	lastId++
	return lastId
}

func New() *EventRepository{
	return &EventRepository{
		eventsList: make(map[uint]*event.Event, startCapacity),
	}
}

func (r *EventRepository) Create(event *event.Event) error {
	event.Id = nextId()
	r.eventsList[event.Id] = event
	return nil
}

func (r *EventRepository) Read(eventId uint) (event *event.Event, err error) {
	event, ok := r.eventsList[eventId]
	if !ok {
		err = calendar_error.NotFound
	}
	return event, err
}


func (r *EventRepository) ReadAll() (map[uint]*event.Event, error) {
	return r.eventsList, nil
}


func (r *EventRepository) Update(event *event.Event) (err error) {

	if savedEvent, ok := r.eventsList[(*event).Id]; ok {
		savedEvent = event
		//r.eventsList[(*event).Id] = event
		_ = savedEvent
	} else {
		err = calendar_error.NotFound
	}
	return err
}


func (r *EventRepository) Delete(eventId uint) (err error) {
	_, ok := r.eventsList[eventId]
	if !ok {
		err = calendar_error.NotFound
	} else {
		delete(r.eventsList, eventId)
	}
	return err
}




