package user

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// TableName const
const TableName = "user"

// User is the user entity
type User struct {
	ID				uint		`gorm:"PRIMARY_KEY" json:"id"`
	Name			string		`gorm:"type:varchar(100);UNIQUE;INDEX" json:"username"`
	Passhash		string		`gorm:"type:bytea" json:"-"`
	CreatedAt		time.Time
	UpdatedAt		time.Time
	DeletedAt		*time.Time	`gorm:"INDEX"`
}

// TableName returns a name of a table
func (e User) TableName() string {
	return TableName
}

// New func is a constructor for the User
func New() *User {
	return &User{}
}

// Validate function for a validation
func (e User) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Name, validation.Required, validation.Length(2, 100), is.Alpha),
		)
}
