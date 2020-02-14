package event

import(
	"strconv"
)

type IEvent interface {

}

type Date struct {
	Year  uint8
	Month uint8
	Day   uint8
}

func (d Date) String() string {
	return strconv.FormatUint(uint64(d.Day), 10) + "." + strconv.FormatUint(uint64(d.Month), 10) + "." + strconv.FormatUint(uint64(d.Year), 10)
}

type Event struct {
	Id		uint
	Name	string
	Date	Date
}

func (e Event) String() string {
	return "#" + strconv.FormatUint(uint64(e.Id), 10) + " " + e.Name + "(" + e.Date.String() + ")"
}
