package inmemory

import(
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/calendarerror"
)

// startCapacity of the repository
const startCapacity = 100

// EventRepository is the event repository :)
type EventRepository struct {
	eventsList map[uint]*event.Event
}

// lastID for entities
var lastID uint

// nextID for entities
func nextID() uint {
	lastID++
	return lastID
}

// New func is a constructor for the EventRepository
func New() *EventRepository{
	return &EventRepository{
		eventsList: make(map[uint]*event.Event, startCapacity),
	}
}

// Create an entity
func (r *EventRepository) Create(event *event.Event) error {

	if !r.IsTimeNotBusy(event) {
		return calendarerror.ErrTimeIsBusy
	}

	event.ID = nextID()
	r.eventsList[event.ID] = event
	return nil
}

// Read an entity
func (r *EventRepository) Read(eventID uint) (event *event.Event, err error) {
	event, ok := r.eventsList[eventID]
	if !ok {
		err = calendarerror.ErrNotFound
	}
	return event, err
}

// ReadAll entities
func (r *EventRepository) ReadAll() (map[uint]*event.Event, error) {
	return r.eventsList, nil
}

// Update an entity
func (r *EventRepository) Update(event *event.Event) (err error) {

	if _, ok := r.eventsList[(*event).ID]; ok {

		if !r.IsTimeNotBusy(event) {
			return calendarerror.ErrTimeIsBusy
		}

		r.eventsList[(*event).ID] = event
	} else {
		err = calendarerror.ErrNotFound
	}
	return err
}

// Delete an entity
func (r *EventRepository) Delete(eventID uint) (err error) {
	_, ok := r.eventsList[eventID]
	if !ok {
		err = calendarerror.ErrNotFound
	} else {
		delete(r.eventsList, eventID)
	}
	return err
}

// IsTimeNotBusy func is check if time for an event is not alredy busy
func (r *EventRepository) IsTimeNotBusy(event *event.Event) bool {
	eventEnd := event.Time.Add(event.Duration)

	for _, e := range r.eventsList {
		if e.ID == event.ID {
			continue
		}
		end := e.Time.Add(e.Duration)

		if (event.Time.After(e.Time) && event.Time.Before(end)) || (event.Time.Equal(e.Time)) || (eventEnd.After(e.Time) && eventEnd.Before(end)) || (eventEnd.Equal(end)) || (event.Time.Before(e.Time) && eventEnd.After(end)) {
			return false
		}
	}
	return true
}