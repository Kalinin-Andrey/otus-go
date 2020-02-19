package in_memory

import (
	"reflect"
	"testing"
	"time"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/calendar_error"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
)

var repo event.IEventRepository

var oneHour, _ = time.ParseDuration("1h")
var twoHour, _ = time.ParseDuration("2h")

var events = map[uint]*event.Event{
	1: &event.Event{
		Id: 1,
		Name: "Event1",
		Date: time.Date(2020, 01, 01, 01, 0, 0, 0, nil),
		Duration:	oneHour,
	},
	2: &event.Event{
		Id: 2,
		Name: "Event2",
		Date: time.Date(2020, 02, 02, 02, 0, 0, 0, nil),
		Duration:	oneHour,
	},
	3: &event.Event{
		Id: 3,
		Name: "Event3",
		Date: time.Date(2020, 03, 03, 03, 0, 0, 0, nil),
		Duration:	oneHour,
	},
	4: &event.Event{
		Id: 4,
		Name: "Event4",
		Date: time.Date(2020, 04, 04, 04, 0, 0, 0, nil),
		Duration:	oneHour,
	},
	5: &event.Event{
		Id: 5,
		Name: "Event5",
		Date: time.Date(2020, 05, 05, 05, 0, 0, 0, nil),
		Duration:	oneHour,
	},
}

func init() {
	repo = New()
}

func TestEventRepository_Create(t *testing.T) {

	for _, event := range events {
		err := repo.Create(event)
		if err != nil {
			t.Fatalf("func Create has an error: %s", err)
		}
	}
}

func TestEventRepository_Read(t *testing.T) {

	for i, event := range events {
		savedEvent, err := repo.Read(i)
		if err != nil {
			t.Fatalf("func Read has an error: %s", err)
		}
		if ! reflect.DeepEqual(event, *savedEvent) {
			t.Errorf("Saved event is not equal, expected %v, has given %v", event, *savedEvent)
		}
	}
}

func TestEventRepository_ReadAll(t *testing.T) {
	savedEvents, err := repo.ReadAll()
	if err != nil {
		t.Fatalf("func ReadAll has an error: %s", err)
	}
	if ! reflect.DeepEqual(events, savedEvents) {
		t.Errorf("Saved event is not equal, expected %v, has given %v", events, savedEvents)
	}
}

func TestEventRepository_Update(t *testing.T) {
	newEvent := events[2]
	newEvent.Id = 1
	err := repo.Update(newEvent)
	if err != nil {
		t.Fatalf("func Update has an error: %s", err)
	}
	savedNewEvent, err := repo.Read(newEvent.Id)
	if err != nil {
		t.Fatalf("func Read has an error: %s", err)
	}

	if ! reflect.DeepEqual(newEvent, savedNewEvent) {
		t.Errorf("Saved event is not equal, expected %v, has given %v", newEvent, savedNewEvent)
	}
}


func TestEventRepository_Delete(t *testing.T) {
	eventId := uint(1)
	err := repo.Delete(eventId)
	if err != nil {
		t.Fatalf("func Delete has an error: %s", err)
	}
	_, err = repo.Read(eventId)
	if err != calendar_error.NotFound {
		t.Errorf("Read deleted event has return unexpected error, expected error %v, has given %v", calendar_error.NotFound, err)
	}
}

