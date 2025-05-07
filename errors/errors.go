// dui-go/errors/errors.go
//
// Responsibilities of this file:
//   - Define ErrorDetail and APIError types.
//   - Provide functions and methods to create, compare, and augment errors.
//   - Offer predefined common error values for standard HTTP status codes.

package errors

import (
	"fmt"
	"strings"
	"sync"
)

// ErrorDetail represents additional context about the error.
// Reason gives a short machine-friendly key for the cause,
// and Message provides a human-readable explanation.
type ErrorDetail struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// APIError represents an error with a status code and an associated message.
// It can hold multiple ErrorDetail entries for in-depth debugging info.
// Modifications to Details are protected by a mutex.
type APIError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
	mu      sync.Mutex
}

// New creates a new APIError. If message is empty, it defaults to "unknown error".
// Optionally, initial details can be provided.
func New(code int, message string, details ...ErrorDetail) *APIError {
	if message == "" {
		message = "unknown error"
	}
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Error returns a formatted string with code, message, and all detail entries.
func (e *APIError) Error() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	detailMessages := make([]string, len(e.Details))
	for i, detail := range e.Details {
		detailMessages[i] = fmt.Sprintf("%s: %s", detail.Reason, detail.Message)
	}

	if len(detailMessages) > 0 {
		return fmt.Sprintf("APIError: code=%d, message=%q, details=[%s]",
			e.Code, e.Message, strings.Join(detailMessages, "; "))
	}
	return fmt.Sprintf("APIError: code=%d, message=%q", e.Code, e.Message)
}

// Is returns true if target is an APIError with the same code and message.
func (e *APIError) Is(target error) bool {
	apiErr, ok := target.(*APIError)
	return ok && apiErr.Code == e.Code && apiErr.Message == e.Message
}

// WithDetails appends details to the APIError and returns it. Thread-safe.
func (e *APIError) WithDetails(details ...ErrorDetail) *APIError {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Details = append(e.Details, details...)
	return e
}

// Predefined common errors.
var (
	ErrBadRequest   = New(400, "bad request")
	ErrUnauthorized = New(401, "unauthorized")
	ErrForbidden    = New(403, "forbidden")
	ErrNotFound     = New(404, "resource not found")
	ErrServerError  = New(500, "internal server error")
)
