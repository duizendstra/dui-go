# dui-go

`dui-go` is a foundational library of utilities designed to streamline the development of Go-based services, particularly within the context of the Google Cloud Platform. It provides essential components for common tasks such as authentication, caching, environment variable management, structured error handling, key-value data storage, and logging.

## Packages

The `dui-go` library is organized into several packages, each addressing a specific area of functionality:

### Authentication (`github.com/duizendstra/dui-go/authentication`)

This package provides a robust mechanism for managing authentication tokens. It defines interfaces for `Token`, `TokenManager`, and `TokenFetcher`, enabling flexible token retrieval, caching, and refresh strategies.

**Key Features:**

*   **`Token` struct:** Represents a token with its value and expiration time.
*   **`TokenManager` interface:** Defines methods for retrieving and managing tokens.
*   **`TokenFetcher` interface:** Abstracts the process of fetching tokens from an external source.
*   **Thread-safe `TokenManager` implementation:** Ensures safe concurrent access to tokens.

### Cache (`github.com/duizendstra/dui-go/cache`)

This package offers a generic `Cache` interface for in-memory and external caching. It provides a basic in-memory cache implementation.

**Key Features:**

*   **`Cache` interface:** Defines basic `Get`, `Set`, `SetAll`, `GetAll`, and `Flush` operations for caching.
*   **`InMemoryCache`:** A thread-safe in-memory cache implementation.

### Env (`github.com/duizendstra/dui-go/env`)

This package simplifies loading environment variables into Go structs in a type-safe manner. It uses reflection and struct tags for configuration, avoiding external dependencies.

**Key Features:**

*   **`Process(spec interface{})` function:** Populates the fields of a struct pointer (`spec`) based on environment variables.
*   **Tag-Based Configuration:**
    *   `env:"VAR_NAME"`: Specifies the environment variable name (defaults to uppercase field name).
    *   `envDefault:"value"`: Provides a default value if the environment variable is not set or empty.
    *   `envRequired:"true"`: Marks the variable as mandatory; returns an error if missing or empty.
*   **Type Safety:** Loads values directly into struct fields with their defined types.
*   **Supported Types:** Currently handles `string`, `int`, `int64`, and `bool`.
*   **Error Reporting:** Returns errors for missing required variables, parsing issues, or unsupported types.

**Example Usage:**

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/duizendstra/dui-go/env"
)

type Config struct {
	Host      string `env:"APP_HOST" envDefault:"localhost"`
	Port      int    `env:"APP_PORT" envRequired:"true"`
	DebugMode bool   `env:"DEBUG_MODE" envDefault:"false"`
}

func main() {
	// Example: Set a required environment variable
	os.Setenv("APP_PORT", "8080")
	defer os.Unsetenv("APP_PORT")

	var cfg Config
	err := env.Process(&cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Host: %s, Port: %d, Debug: %t\n", cfg.Host, cfg.Port, cfg.DebugMode)
	// Output: Host: localhost, Port: 8080, Debug: false
}
```

### Errors (`github.com/duizendstra/dui-go/errors`)

This package enhances Go's error handling by introducing structured error types with codes and details, suitable for API responses.

**Key Features:**

*   **`APIError` struct:** Represents an error with a code, message, and a slice of `ErrorDetail` structs (each with a reason and message).
*   **`New()` function:** Creates new `APIError` instances.
*   **`WithDetails()` method:** Allows appending additional details to an existing `APIError`.
*   Predefined common errors (e.g., `ErrBadRequest`, `ErrNotFound`).

### Firestore (`github.com/duizendstra/dui-go/firestore`)

This package provides a key-value store abstraction (`KV` interface) over Google Cloud Firestore. It enables storing and retrieving string values. The underlying Google Cloud Firestore client library may provide its own retry mechanisms for transient errors.

**Key Features:**

*   **`KV` interface:** Defines `Get`, `Set`, and `Close` operations for a key-value store.
*   **`FirestoreKV`:** Implements the `KV` interface using Google Cloud Firestore. It stores string values in a specified collection, where each document contains a "value" field holding the string.
*   If a key does not exist, `Get` returns an empty string and no error.

### Store (`github.com/duizendstra/dui-go/store`)

This package defines a generic key-value `Store` interface and provides a `FirestoreStore` implementation that uses the `firestore.KV` interface as its backend.

**Key Features:**

*   **`Store` interface:** Defines basic `Get`, `Set`, and `Close` operations.
*   **`FirestoreStore`:** Implements `Store` using a `firestore.KV` compatible backend.

### Logging/Cloudlogging (`github.com/duizendstra/dui-go/logging/cloudlogging`)

This package provides a logging handler specifically designed for Google Cloud Logging, integrating with `log/slog`.

**Key Features:**

*   **`NewCloudLoggingHandler()`:** Creates a `slog.Handler` that formats log entries as JSON compatible with Cloud Logging.
*   **Structured logging:** Maps standard `slog` fields and levels to Cloud Logging's expected format (e.g., "severity", "message").
*   **Trace Integration:** Includes `WithCloudTraceContext` HTTP middleware to extract `X-Cloud-Trace-Context` headers and inject trace information into log entries.
*   Custom log levels (e.g., `LevelNotice`, `LevelCritical`) aligning with Cloud Logging severities.

## Getting Started

1.  **Installation:**

    ```bash
    go get github.com/duizendstra/dui-go
    ```

2.  **Usage:**

    Import the desired packages into your Go code:

    ```go
    import (
        "github.com/duizendstra/dui-go/authentication"
        "github.com/duizendstra/dui-go/cache"
        "github.com/duizendstra/dui-go/env"
        "github.com/duizendstra/dui-go/errors"
        "github.com/duizendstra/dui-go/firestore" // Direct use or via store
        "github.com/duizendstra/dui-go/logging/cloudlogging"
        "github.com/duizendstra/dui-go/store"
    )
    ```

    Refer to the individual package documentation (`doc.go` files) and examples (`example_test.go` files) for detailed usage instructions.

## Building and Testing

To build the `dui-go` library, run:

```bash
go build ./...
```

To test the library:

```bash
go test ./...
```

To run tests with the race detector:
```bash
go test -race ./...
```

## Code of Conduct

Act professionally and respectfully. Be kind, considerate, and welcoming. Harassment or exclusionary behavior will not be tolerated.

## Related Files

*   `CHANGELOG.md`: Tracks notable changes for each release.
*   `CONTRIBUTING.md`: Contribution guidelines.
*   `ROADMAP.md`: Future plans for the library.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
