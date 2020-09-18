package web

import (
	"github.com/pkg/errors"
	"net/http"
)

// GenericErrorResponse is the API response form for generic failures.
type GenericErrorResponse struct {
	Error string `json:"error"`
}

// GenericError is used to pass a request error during the request with specific context.
type GenericError struct {
	Err        error
	StatusCode int
}

func (err *GenericError) Error() string {
	return err.Err.Error()
}

// NewGenericError wraps a provided error with an HTTP status code.
// It should be used for expected errors.
func NewGenericError(err error, status int) error {
	return &GenericError{err, status}
}

// FieldError is used to indicate an error with a specific request field like
// if a field is mandatory and is not passed or if a field is passed but with invalid type.
type FieldError struct {
	Error string `json:"error"`
	Field string `json:"field"`
}

// FieldValidationErrorResponse is the API response form for field validation failures.
type FieldValidationErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields"`
}

// FieldsValidationError is used to pass a request error during the request with specific context.
type FieldsValidationError struct {
	GenericError
	Fields []FieldError
}

// FieldsValidationError extends a GenericError with field errors and set a 400 BAD REQUEST status
func NewFieldsValidationError(fields []FieldError) error {
	return &FieldsValidationError{
		GenericError: GenericError{errors.New("field validation error"), http.StatusBadRequest},
		Fields:       fields,
	}
}

// RedirectError is used to redirect a request instead of returning a response.
type RedirectError struct {
	Url        string
	Err        error
	StatusCode int
}

func (err *RedirectError) Error() string {
	return err.Err.Error()
}

// NewGenericError wraps a provided error with an HTTP status code.
// It should be used for expected errors.
func NewRedirectError(url string, status int) error {
	return &RedirectError{url, errors.New(""), status}
}

// shutdown is a type used to help with the graceful termination of the service.
type shutdown struct {
	Message string
}

// NewShutdownError returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdown{message}
}

func (s *shutdown) Error() string {
	return s.Message
}

// IsShutdown checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}
	return false
}
