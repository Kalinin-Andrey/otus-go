package grpcerror

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/calendarpb"
)

// StatusOK is status for case when every thing is OK
var StatusOK = &calendarpb.Status{
	OK:	true,
}

// StatusBadRequest is status for case when bad request
var StatusBadRequest = &calendarpb.Status{
	OK:	false,
	Error: "bad request",
}

// StatusNotFound is status for case when entity not found
var StatusNotFound = &calendarpb.Status{
	OK:	false,
	Error: "not found",
}

// StatusInternalServerError is status for case when smth went wrong
var StatusInternalServerError = &calendarpb.Status{
	OK:	false,
	Error: "internal server error",
}

// ErrBadRequestResponseEvents is error for case when bad request
var ErrBadRequestResponseEvents = &calendarpb.ResponseEvents{
	Status: StatusBadRequest,
}

// ErrBadRequestResponseEvent is error for case when bad request
var ErrBadRequestResponseEvent = &calendarpb.ResponseEvent{
	Status: StatusBadRequest,
}

// ErrNotFoundResponseEvents is error for case when entity not found
var ErrNotFoundResponseEvents = &calendarpb.ResponseEvents{
	Status: StatusNotFound,
}

// ErrNotFoundResponseEvent is error for case when entity not found
var ErrNotFoundResponseEvent = &calendarpb.ResponseEvent{
	Status: StatusNotFound,
}

// ErrInternalResponseEvents is error for case when smth went wrong
var ErrInternalResponseEvents = &calendarpb.ResponseEvents{
	Status: StatusInternalServerError,
}

// ErrInternalResponseEvent is error for case when smth went wrong
var ErrInternalResponseEvent = &calendarpb.ResponseEvent{
	Status: StatusInternalServerError,
}
