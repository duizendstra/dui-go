// dui-go/errors/example_test.go
//
// This example demonstrates how to create and use an APIError with details.
// It does not produce output, as requested.
//
// In real usage, you might return this error from an API handler
// or convert it into an HTTP response.

package errors

import "fmt"

func ExampleAPIError() {
	// Create a new error with a code and message
	err := New(401, "unauthorized")

	// Add details to provide more context
	err = err.WithDetails(ErrorDetail{Reason: "AUTH_FAILED", Message: "Invalid token"})

	// In a real scenario, you might log this error, return it in an HTTP response,
	// or wrap it with additional context.
	_ = fmt.Sprintf("Handle error: %v", err)

	// No output is produced by design.
}
