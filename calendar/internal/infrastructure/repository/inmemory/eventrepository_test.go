package inmemory

import (
	calendarerror2 "github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/calendarerror"
	"reflect"
	"testing"
	"time"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
)

var repo event.IEventRepository

var halfHour time.Duration
var oneHour time.Duration
var twoHour time.Duration
var threeHour time.Duration
var locationMoscow *time.Location

var events map[uint]*event.Event
var eventsIDs []uint

func init() {
	repo = New()
	var err error

	halfHour, err = time.ParseDuration("0.5h")
	if err != nil {
		panic(err)
	}

	oneHour, err = time.ParseDuration("1h")
	if err != nil {
		panic(err)
	}

	twoHour, err = time.ParseDuration("2h")
	if err != nil {
		panic(err)
	}

	threeHour, err = time.ParseDuration("3h")
	if err != nil {
		panic(err)
	}

	locationMoscow, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	events = map[uint]*event.Event{
		1: &event.Event{
			ID:       1,
			Name:     "Event1",
			Time:     time.Date(2020, 01, 01, 01, 0, 0, 0, locationMoscow),
			Duration: oneHour,
		},
		2: &event.Event{
			ID:       2,
			Name:     "Event2",
			Time:     time.Date(2020, 02, 02, 02, 0, 0, 0, locationMoscow),
			Duration: oneHour,
		},
		3: &event.Event{
			ID:       3,
			Name:     "Event3",
			Time:     time.Date(2020, 03, 03, 03, 0, 0, 0, locationMoscow),
			Duration: oneHour,
		},
		4: &event.Event{
			ID:       4,
			Name:     "Event4",
			Time:     time.Date(2020, 04, 04, 04, 0, 0, 0, locationMoscow),
			Duration: oneHour,
		},
		5: &event.Event{
			ID:       5,
			Name:     "Event5",
			Time:     time.Date(2020, 05, 05, 05, 0, 0, 0, locationMoscow),
			Duration: oneHour,
		},
	}
	eventsIDs = []uint{1, 2, 3, 4, 5}
}

func TestEventRepository_Create(t *testing.T) {

	for _, eventID := range eventsIDs {
		_, err := repo.Create(events[eventID])
		if err != nil {
			t.Fatalf("func Create has an error: %s", err)
		}
	}
}

func TestEventRepository_Read(t *testing.T) {

	for _, event := range events {
		savedEvent, err := repo.Read(event.ID)
		if err != nil {
			t.Fatalf("func Read has an error: %s", err)
		}
		if !reflect.DeepEqual(*event, *savedEvent) {
			t.Errorf("Saved event is not equal, expected %v, has given %v", *event, *savedEvent)
		}
	}
}

func TestEventRepository_ReadAll(t *testing.T) {
	savedEvents, err := repo.ReadAll()
	if err != nil {
		t.Fatalf("func ReadAll has an error: %s", err)
	}
	if !reflect.DeepEqual(events, savedEvents) {
		t.Errorf("Saved event is not equal, expected %v, has given %v", events, savedEvents)
	}
}

func TestEventRepository_Update(t *testing.T) {
	updatedEvent := *events[2]
	updatedEvent.Name = "Updated" + updatedEvent.Name
	err := repo.Update(&updatedEvent)
	if err != nil {
		t.Fatalf("func Update has an error: %s", err)
	}
	savedEvent, err := repo.Read(updatedEvent.ID)
	if err != nil {
		t.Fatalf("func Read has an error: %s", err)
	}

	if !reflect.DeepEqual(updatedEvent, *savedEvent) {
		t.Errorf("Saved event is not equal, expected %v, has given %v", updatedEvent, *savedEvent)
	}
}

func TestEventRepository_Delete(t *testing.T) {
	eventID := uint(3)
	err := repo.Delete(eventID)
	if err != nil {
		t.Fatalf("func Delete has an error: %s", err)
	}
	_, err = repo.Read(eventID)
	if err != calendarerror2.ErrNotFound {
		t.Errorf("Read deleted event has return unexpected error, expected error %v, has given %v", calendarerror2.ErrNotFound, err)
	}
}

func TestEventRepository_TimeIsBusy(t *testing.T) {
	event := events[1]
	updatedEvent := events[2]

	updatedEvent.Time = event.Time
	err := repo.Update(updatedEvent)
	if err != calendarerror2.ErrTimeIsBusy {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", calendarerror2.ErrTimeIsBusy, err)
	}

	updatedEvent.Time = event.Time.Add(halfHour)
	err = repo.Update(updatedEvent)
	if err != calendarerror2.ErrTimeIsBusy {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", calendarerror2.ErrTimeIsBusy, err)
	}

	updatedEvent.Time = time.Date(2020, 01, 01, 0, 30, 0, 0, locationMoscow)
	err = repo.Update(updatedEvent)
	if err != calendarerror2.ErrTimeIsBusy {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", calendarerror2.ErrTimeIsBusy, err)
	}

	updatedEvent.Time = time.Date(2020, 01, 01, 00, 0, 0, 0, locationMoscow)
	updatedEvent.Duration = twoHour
	err = repo.Update(updatedEvent)
	if err != calendarerror2.ErrTimeIsBusy {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", calendarerror2.ErrTimeIsBusy, err)
	}

	updatedEvent.Time = time.Date(2020, 01, 01, 00, 0, 0, 0, locationMoscow)
	updatedEvent.Duration = threeHour
	err = repo.Update(updatedEvent)
	if err != calendarerror2.ErrTimeIsBusy {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", calendarerror2.ErrTimeIsBusy, err)
	}

	updatedEvent.Time = time.Date(2020, 01, 01, 00, 0, 0, 0, locationMoscow)
	updatedEvent.Duration = oneHour
	err = repo.Update(updatedEvent)
	if err != nil {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", nil, err)
	}

	updatedEvent.Time = time.Date(2020, 01, 01, 20, 0, 0, 0, locationMoscow)
	err = repo.Update(updatedEvent)
	if err != nil {
		t.Errorf("Update event has return unexpected error, expected error %v, has given %v", nil, err)
	}

}



