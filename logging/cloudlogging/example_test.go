// go-dui/pkg/logging/cloudlogging/example_test.go
//
// This file provides a runnable example demonstrating how to use the
// CloudLoggingHandler and WithCloudTraceContext middleware together.

package cloudlogging

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
)

func ExampleNewCloudLoggingHandler() {
	// Create a CloudLoggingHandler that outputs JSON logs recognized by Cloud Logging.
	handler := NewCloudLoggingHandler("example-component")

	// Create a logger using our handler.
	logger := slog.New(handler)

	// Wrap a simple HTTP handler with the WithCloudTraceContext middleware.
	// In a real application, this would be attached to your HTTP server.
	wrappedHandler := WithCloudTraceContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the logger to write logs. This will include trace info if present.
		logger.InfoContext(r.Context(), "Received request", "path", r.URL.Path)

		if r.URL.Path != "/expected-path" {
			// If the path doesn't match, log and return an error.
			logger.ErrorContext(r.Context(), "Unexpected path accessed", "path", r.URL.Path)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Respond successfully and log the success.
		fmt.Fprintln(w, "Hello, Cloud Logging with Trace!")
		logger.InfoContext(r.Context(), "Request handled successfully")
	}))

	// Create a simulated request with trace headers.
	// Typically, these headers are set by Google Cloud's load balancer or App Engine.
	testReq := httptest.NewRequest("GET", "http://example.com/expected-path", nil)
	testReq.Header.Set("X-Cloud-Trace-Context", "abcdef1234567890abcdef1234567890/123;o=1")

	// Serve the request using httptest.
	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, testReq)

	// Inspect the response.
	responseBody := rr.Body.String()
	fmt.Println("Response body:", responseBody)

	// Log outside a request context to show how logs differ without trace info.
	logger.InfoContext(context.Background(), "Logging outside of a request context")

	// Note: The actual logs will be printed to stderr in JSON format.
	// The trace info will be included in the logs related to the handled request.
}
