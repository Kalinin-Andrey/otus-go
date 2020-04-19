package event

import (
	"strconv"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

)

// Event entity
type Event struct {
	ID				uint					`gorm:"PRIMARY_KEY" json:"id"`
	UserID			uint					`sql:"index" json:"userId"`
	Title			string					`gorm:"type:varchar(100)" json:"title"`
	Description		*string					`json:"description,omitempty"`
	Time			time.Time				`sql:"index"`
	Duration		time.Duration
	NoticePeriod	*time.Duration			`json:"noticePeriod,omitempty"`

	CreatedAt		time.Time
	UpdatedAt		time.Time
	DeletedAt		*time.Time				`sql:"index"`
}

const (
	// TableName const
	TableName	= "event"
)

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
