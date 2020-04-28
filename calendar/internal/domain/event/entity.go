package event

import (
	"encoding/json"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/notification"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Event entity
type Event struct {
	ID				uint					`json:"id"`
	UserID			uint					`db:"user_id" json:"userId"`
	Title			string					`json:"title"`
	Description		*string					`json:"description,omitempty"`
	Time			time.Time
	Duration		time.Duration			`json:"duration"`
	NoticePeriod	*time.Duration			`db:"notice_period" json:"noticePeriod,omitempty"`
}

const (
	// TableName const
	TableName	= "event"
)

// QueryCondition struct for defining a query condition
type QueryCondition struct {
	Where	*WhereCondition
}
// WhereCondition struct
type WhereCondition struct {
	Time	*WhereConditionTime
}
// WhereConditionTime struct
type WhereConditionTime struct {
	Between	*[2]time.Time
}

// Validate func
func (e Event) Validate() error {

	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Required, validation.Length(2, 100), is.PrintableASCII),
	)
}

// TableName func
func (e Event) TableName() string {
	return TableName
}

// New func is a constructor
func New() *Event {
	return &Event{}
}

// String func is a func for the Stringer interface
func (e Event) String() string {
	return "#" + strconv.FormatUint(uint64(e.ID), 10) + " " + e.Title + "(" + e.Time.String() + ")"
}
// UnmarshalJSON is unmarshal func for an entity
func (e *Event) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		ID				uint
		UserID			uint
		Title			string
		Description		*string
		Time			string
		Duration		string
		NoticePeriod	*string
	}
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	e.ID		= tmp.ID
	e.UserID	= tmp.UserID
	e.Title		= tmp.Title
	e.Time, err	= time.Parse(time.RFC3339, tmp.Time)
	if err != nil {
		return errors.Wrapf(err, "Can not parse time for a field Time: %v", tmp.Time)
	}
	duration, err := time.ParseDuration(tmp.Duration)
	if err != nil {
		return errors.Wrapf(err, "Can not parse duration for a field Duration: %v", tmp.Duration)
	}
	e.Duration	= duration

	if tmp.Description != nil {
		e.Description = tmp.Description
	}

	if tmp.NoticePeriod != nil {
		noticePeriod, err := time.ParseDuration(*tmp.NoticePeriod)
		if err != nil {
			return errors.Wrapf(err, "Can not parse duration for a field NoticePeriod: %v", tmp.NoticePeriod)
		}
		e.NoticePeriod	= &noticePeriod
	}

	return err
}

// Notification returns a Notification object for a current Event object
func (e Event) Notification() *notification.Notification {
	return &notification.Notification{
		EventID:	e.ID,
		UserID:		e.UserID,
		Title:		e.Title,
		Time:		e.Time,
	}
}

