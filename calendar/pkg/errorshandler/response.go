package errorshandler

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"sort"
)

// Response is the response that represents an error.
type Response struct {
	Status  int         `json:"-"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success returns the struct for success response
func Success() Response {
	return Response{
		Status: http.StatusOK,
		Message: "success",
	}
}

// SuccessMessage returns the struct for success response with only message
func SuccessMessage() Response {
	return Response{
		Message: "success",
	}
}

// Error is required by the error interface.
func (e Response) Error() string {
	return e.Message
}

// StatusCode is required by routing.HTTPError interface.
func (e Response) StatusCode() int {
	return e.Status
}

// InternalServerError creates a new error response representing an internal server error (HTTP 500)
func InternalServerError(msg string) Response {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return Response{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

// NotFound creates a new error response representing a resource-not-found error (HTTP 404)
func NotFound(msg string) Response {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return Response{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

// Unauthorized creates a new error response representing an authentication/authorization failure (HTTP 401)
func Unauthorized(msg string) Response {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return Response{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}

// Forbidden creates a new error response representing an authorization failure (HTTP 403)
func Forbidden(msg string) Response {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return Response{
		Status:  http.StatusForbidden,
		Message: msg,
	}
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest(msg string) Response {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return Response{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

type invalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// InvalidInput creates a new error response representing a data validation error (HTTP 400).
func InvalidInput(errs validation.Errors) Response {
	var details []invalidField
	var fields []string
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		details = append(details, invalidField{
			Field: field,
			Error: errs[field].Error(),
		})
	}

	return Response{
		Status:  http.StatusBadRequest,
		Message: "There is some problem with the data you submitted.",
		Details: details,
	}
}
