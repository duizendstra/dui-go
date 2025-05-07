// dui-go/errors/doc.go
//
// Package errors provides a structured error type (APIError) with HTTP-like
// status codes and optional detailed context. This makes it easier for APIs
// to produce consistent, informative error responses.
//
// Usage:
//   err := errors.New(400, "bad request")
//   err = err.WithDetails(errors.ErrorDetail{Reason: "INVALID_INPUT", Message: "Missing 'id' parameter"})
//
// APIError can be compared using errors.Is, and details can be appended as needed.
// Predefined common errors (ErrBadRequest, ErrNotFound, etc.) are also provided.
//
// This package is thread-safe for adding details and reading error information.
// Its error messages and details are suitable for logging, or returning to
// clients in JSON format.

package errors
