package event

import(
	"strconv"
	"time"
)

//type IEvent interface {}

// Event is a entity of a event
type Event struct {
	ID       uint
	Name     string
	Time     time.Time
	Duration time.Duration
}

// New func is a constructor for the Event
func New() *Event{
	return &Event{}
}

// String func is a func for the Stringer interface
func (e Event) String() string {
	return "#" + strconv.FormatUint(uint64(e.ID), 10) + " " + e.Name + "(" + e.Time.String() + ")"
}
