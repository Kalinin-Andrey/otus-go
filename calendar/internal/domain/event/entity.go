package event

import(
	"strconv"
	"time"
)

//type IEvent interface {}

type Event struct {
	Id			uint
	Name		string
	Date		time.Time
	Duration	time.Duration
}

func New() *Event{
	return &Event{}
}

func (e Event) String() string {
	return "#" + strconv.FormatUint(uint64(e.Id), 10) + " " + e.Name + "(" + e.Date.String() + ")"
}
