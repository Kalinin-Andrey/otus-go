package notification

import (
	"time"
)

// Notification value object
type Notification struct {
	EventID			uint					`json:"eventId"`
	UserID			uint					`json:"userId"`
	Title			string					`json:"title"`
	Time			time.Time				`sql:"index"`
}

// New func is a constructor
func New() *Notification {
	return &Notification{}
}

// String func is a func for the Stringer interface
func (e Notification) String() string {
	return e.Title + " (" + e.Time.String() + ")"
}
