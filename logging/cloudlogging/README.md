# Go Middleware for Google Cloud Logging & Trace

[![Go Report Card](https://goreportcard.com/badge/github.com/your-repo/dui-go/logging/cloudlogging)](https://goreportcard.com/report/github.com/your-repo/dui-go/logging/cloudlogging)

This Go package provides middleware and a `slog.Handler` to seamlessly integrate standard Go applications with Google Cloud Logging and Google Cloud Trace. It is designed to be a lightweight, production-ready solution for applications running in Google Cloud environments like Cloud Run, GKE, and App Engine.

The primary goal is to produce structured JSON logs that are correctly interpreted by Cloud Logging, including proper severity levels and automatic correlation with trace data.

## Features

-   **`CloudLoggingHandler`**: A `slog.Handler` that formats log entries into the JSON structure expected by Google Cloud Logging.
-   **Automatic Severity Mapping**: Translates standard `slog` levels (and custom ones like `NOTICE`, `CRITICAL`) to the correct `severity` field for Cloud Logging.
-   **Trace Correlation**: Automatically injects trace and span IDs into logs, linking them to the parent request trace.
-   **`WithCloudTraceContext` Middleware**: An HTTP middleware that parses the `X-Cloud-Trace-Context` header from incoming requests and makes trace data available in the `context.Context`.
-   **Source Location**: Automatically includes the file, line, and function name for every log entry under the `logging.googleapis.com/sourceLocation` field.
-   **Environment-Based Configuration**: Configure the log level dynamically using the `LOG_LEVEL` environment variable.
-   **Cloud Run Jobs Support**: Provides the `WithTrace` helper function to easily correlate logs in non-HTTP environments where the trace context must be manually initiated.
-   **Testable**: Includes public helpers (`NewCloudLoggingHandlerForTest`, `SetProjectIDFetcher`) to make testing your application's logging behavior straightforward.

## Installation

```sh
go get github.com/your-repo/dui-go/logging/cloudlogging
```

## Quick Start: HTTP Services

For a typical web server running in Google Cloud, wrap your main HTTP handler with the `WithCloudTraceContext` middleware and set the `CloudLoggingHandler` as the default `slog` logger.

```go
package main

import (
	"log/slog"
	"net/http"
	"os"

	"your-repo/path/to/cloudlogging"
)

func main() {
	// 1. Create a logger using the CloudLoggingHandler.
	// The "my-web-app" component name will be added to all log entries.
	logger := slog.New(cloudlogging.NewCloudLoggingHandler("my-web-app"))

	// 2. Set it as the global default logger.
	slog.SetDefault(logger)

	// 3. Define your application's HTTP handler.
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log messages using the default logger. Because of the middleware,
		// the request context contains trace information, which will be
		// automatically added to the log entry.
		slog.InfoContext(r.Context(), "Handling request", "path", r.URL.Path)

		w.Write([]byte("Hello, World!"))

		slog.InfoContext(r.Context(), "Request handled successfully")
	})

	// 4. Wrap your handler with the trace context middleware.
	http.Handle("/", cloudlogging.WithCloudTraceContext(helloHandler))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
```

When a request with the header `X-Cloud-Trace-Context: abcdef123/456;o=1` is received, the logs produced by the handler will look like this in Cloud Logging:

```json
{
  "severity": "INFO",
  "message": "Handling request",
  "timestamp": "...",
  "component": "my-web-app",
  "path": "/",
  "logging.googleapis.com/sourceLocation": {
    "file": "/path/to/main.go",
    "line": 23,
    "function": "main.main.func1"
  },
  "logging.googleapis.com/trace": "projects/your-gcp-project/traces/abcdef123",
  "logging.googleapis.com/spanId": "456",
  "logging.googleapis.com/trace_sampled": true
}
```

## Configuration

### Log Level

The logger's minimum level is controlled by the **`LOG_LEVEL`** environment variable. If unset or set to an invalid value, it defaults to `INFO`. The comparison is case-insensitive.

| `LOG_LEVEL` Value | `slog.Level`   | Cloud Logging Severity |
| ----------------- | -------------- | ---------------------- |
| `DEBUG`           | `LevelDebug`   | `DEBUG`                |
| `INFO`            | `LevelInfo`    | `INFO`                 |
| `NOTICE`          | `LevelNotice`  | `NOTICE`               |
| `WARN` or `WARNING` | `LevelWarn`    | `WARNING`              |
| `ERROR`           | `LevelError`   | `ERROR`                |
| `CRITICAL`        | `LevelCritical`| `CRITICAL`             |
| `ALERT`           | `LevelAlert`   | `ALERT`                |
| `EMERGENCY`       | `LevelEmergency`| `EMERGENCY`            |

### Google Cloud Project ID

The middleware and handler need your Google Cloud Project ID to format the trace string correctly. The ID is determined in the following order of precedence:

1.  **GCP Metadata Server**: Fetched automatically when running on Google Cloud infrastructure.
2.  **`GOOGLE_CLOUD_PROJECT` Environment Variable**: A fallback if the metadata server is unavailable.
3.  **`"unknown-project"`**: A default value if neither of the above can be resolved.

## Advanced Usage: Cloud Run Jobs

Cloud Run Jobs are not triggered by HTTP requests and therefore don't have an incoming `X-Cloud-Trace-Context` header. To correlate all logs from a single job execution, you can manually create a trace ID at the start of the job and inject it into the context using the provided `WithTrace` helper function. This pattern ensures that all logs from the job can be easily filtered and viewed together.

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid" // Requires: go get github.com/google/uuid
	"your-repo/path/to/cloudlogging"
)

func main() {
	// 1. Set up the production-ready logger.
	logger := slog.New(cloudlogging.NewCloudLoggingHandler("my-cloud-run-job"))
	slog.SetDefault(logger)

	// 2. Create a root context for this specific job execution.
	ctx := context.Background()

	// 3. Create a unique trace ID for this execution.
	// You can use the built-in execution ID from Cloud Run or generate a new UUID.
	executionID := os.Getenv("CLOUD_RUN_EXECUTION")
	if executionID == "" {
		executionID = uuid.NewString()
	}

	// 4. Use the WithTrace helper to inject the trace ID into the context.
	ctx = cloudlogging.WithTrace(ctx, executionID)

	// 5. Run the job's business logic, passing the context down.
	slog.InfoContext(ctx, "Job starting.")
	// ... run job logic with ctx ...
	slog.InfoContext(ctx, "Job finished successfully.")
}
```

## Testing

The package is designed for testability. You can test your application's logging output without writing to `stderr` and without relying on the real GCP metadata service.

-   **`NewCloudLoggingHandlerForTest`**: This constructor allows you to direct log output to any `io.Writer`, like a `bytes.Buffer`, for inspection.

-   **`SetProjectIDFetcher`**: This function allows you to replace the default metadata-based project ID fetcher with a mock implementation, ensuring deterministic behavior in your tests.

### Example Test

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"io"

	"your-repo/path/to/cloudlogging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A mock fetcher for testing.
type mockProjectIDFetcher struct {
	id  string
	err error
}
func (m *mockProjectIDFetcher) ProjectID(ctx context.Context) (string, error) {
	return m.id, m.err
}


func TestMyApplicationLogging(t *testing.T) {
	// Mock the project ID fetcher to return a known value.
	cloudlogging.SetProjectIDFetcher(&mockProjectIDFetcher{id: "test-project", err: nil})

	// Create a buffer to capture log output.
	var logBuffer bytes.Buffer

	// Use the public testing constructor to create a handler that writes to our buffer.
	handler := cloudlogging.NewCloudLoggingHandlerForTest(&logBuffer, "test-component")
	logger := slog.New(handler)

	// Simulate a trace context
	ctx := context.Background()
	ctx = cloudlogging.WithTrace(ctx, "my-test-trace")

	logger.InfoContext(ctx, "This is a test log")

	var logOutput map[string]interface{}
	err := json.Unmarshal(logBuffer.Bytes(), &logOutput)
	require.NoError(t, err, "Log output should be valid JSON")

	assert.Equal(t, "INFO", logOutput["severity"])
	assert.Equal(t, "This is a test log", logOutput["message"])
	assert.Equal(t, "test-component", logOutput["component"])
	assert.Contains(t, logOutput["logging.googleapis.com/trace"], "projects/test-project/traces/my-test-trace")
}