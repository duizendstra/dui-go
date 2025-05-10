// dui-go/errors/doc.go
//
// Package errors provides a structured error type (APIError) with HTTP-like
// status codes and optional detailed context. This makes it easier for services
// consuming this library to produce consistent, informative error responses,
// particularly for APIs.
//
// Usage:
//
//	import "github.com/duizendstra/dui-go/errors"
//
//	func handleRequest() error {
//	    // ... some operation fails ...
//	    apiErr := errors.New(errors.ErrBadRequest.Code, "Invalid user ID provided") // Or use errors.ErrBadRequest directly
//	    apiErr = apiErr.WithDetails(
//	        errors.ErrorDetail{Reason: "INVALID_FORMAT", Message: "User ID must be a UUID."},
//	        errors.ErrorDetail{Reason: "VALUE_TOO_SHORT", Message: "User ID minimum length is 36 characters."},
//	    )
//	    return apiErr
//	}
//
// APIError can be compared using errors.Is for the base error type (e.g., errors.Is(err, errors.ErrBadRequest)),
// and details can be appended as needed using WithDetails.
// Predefined common errors (ErrBadRequest, ErrNotFound, etc.) are also provided for convenience.
//
// This package is thread-safe for adding details and reading error information.
// Its error messages and details are suitable for logging or returning to
// clients in JSON format (due to struct tags in APIError and ErrorDetail).
package errors
