# dui-go

[![Go Reference](https://pkg.go.dev/badge/github.com/duizendstra/dui-go.svg)](https://pkg.go.dev/github.com/duizendstra/dui-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/duizendstra/dui-go)](https://goreportcard.com/report/github.com/duizendstra/dui-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`dui-go` is a foundational library of idiomatic Go utilities designed to streamline the development of robust and maintainable services, particularly those interacting with the Google Cloud Platform.

## Explore in Firebase Studio

Want to browse the `dui-go` library code or contribute? You can quickly open this repository in Google's Firebase Studio (formerly Project IDX), a browser-based development environment:

<!-- Open dui-go in Firebase Studio Button -->
<a href="https://studio.firebase.google.com/import?url=https%3A%2F%2Fgithub.com%2Fduizendstra%2Fdui-go">
  <picture>
    <source
      media="(prefers-color-scheme: dark)"
      srcset="https://cdn.firebasestudio.dev/btn/open_dark_32.svg">
    <source
      media="(prefers-color-scheme: light)"
      srcset="https://cdn.firebasestudio.dev/btn/open_light_32.svg">
    <img
      height="32"
      alt="Open dui-go in Firebase Studio"
      src="https://cdn.firebasestudio.dev/btn/open_blue_32.svg">
  </picture>
</a>
<!-- End Button -->

Clicking the button above will create a new workspace in Firebase Studio with a clone of the `dui-go` repository. This is a convenient way to explore the source code, run tests, or start contributing, leveraging the pre-configured Nix environment if one is defined in `.idx/dev.nix` for the library's development.

## Why `dui-go`?

*   **Reduce Boilerplate:** Provides well-tested, reusable components for common tasks, letting you focus on core application logic.
*   **Idiomatic Go:** Encourages best practices in Go development, including clear API design, explicit error handling, and effective concurrency.
*   **Simplified GCP Integration:** Offers convenient abstractions for services like Firestore and Cloud Logging, including trace context propagation.
*   **Type Safety & Minimal Dependencies:** Emphasizes type safety (e.g., in environment variable loading) and strives for minimal external dependencies.
*   **Testability Focus:** Designed with testability in mind, promoting the use of interfaces and providing mocks.

## Packages

The `dui-go` library is organized into several packages, each addressing a specific area of functionality. Refer to individual package documentation (linked) for detailed usage and `example_test.go` files.

### [Authentication](./authentication/) (`github.com/duizendstra/dui-go/authentication`)

Provides a robust mechanism for managing authentication tokens. It defines interfaces for `Token`, `TokenManager`, and `TokenFetcher`, enabling flexible token retrieval, caching, and refresh strategies.

**Key Features:**
*   `Token` struct for token value and expiration.
*   `TokenManager` interface for managing token lifecycle.
*   `TokenFetcher` interface for abstracting token acquisition.
*   Thread-safe `TokenManager` implementation.

### [Cache](./cache/) (`github.com/duizendstra/dui-go/cache`)

Offers a generic `Cache` interface for in-memory and potentially external caching. Includes a basic thread-safe in-memory cache implementation.

**Key Features:**
*   `Cache` interface (`Get`, `Set`, `SetAll`, `GetAll`, `Flush`).
*   `InMemoryCache` implementation.

### [Env](./env/) (`github.com/duizendstra/dui-go/env`)

Simplifies loading environment variables into Go structs in a type-safe manner using reflection and struct tags, avoiding external dependencies for this common task.

**Key Features:**
*   `Process(spec interface{})` function for populating structs.
*   Tag-Based Configuration (`env`, `envDefault`, `envRequired`).
*   Type safety for `string`, `int`, `int64`, `bool`.
*   Clear error reporting for missing or unparsable variables.

*(See a detailed example for the `env` package under "Basic Usage Example" below.)*

### [Errors](./errors/) (`github.com/duizendstra/dui-go/errors`)

Enhances Go's error handling by introducing structured `APIError` types with codes and details, particularly useful for building consistent API error responses.

**Key Features:**
*   `APIError` struct with `Code`, `Message`, and `Details`.
*   `New()` and `WithDetails()` for error creation and augmentation.
*   Predefined common errors (e.g., `ErrBadRequest`, `ErrNotFound`).

### [Firestore](./firestore/) (`github.com/duizendstra/dui-go/firestore`)

Provides a simplified key-value store abstraction (`KV` interface) over Google Cloud Firestore, enabling easy storage and retrieval of string values.

**Key Features:**
*   `KV` interface (`Get`, `Set`, `Close`).
*   `FirestoreKV` implementation using a Firestore collection.
*   Handles non-existent keys by returning an empty string without error.

### [Store](./store/) (`github.com/duizendstra/dui-go/store`)

Defines a generic key-value `Store` interface and provides a `FirestoreStore` implementation that utilizes a `firestore.KV` compatible backend. This promotes abstraction over specific storage mechanisms.

**Key Features:**
*   `Store` interface (`Get`, `Set`, `Close`).
*   `FirestoreStore` implementation.

### [Logging/Cloudlogging](./logging/cloudlogging/) (`github.com/duizendstra/dui-go/logging/cloudlogging`)

Provides a logging handler specifically designed for Google Cloud Logging, integrating seamlessly with Go's standard `log/slog` package.

**Key Features:**
*   `NewCloudLoggingHandler()` for GCP-compatible JSON logs.
*   Structured logging with standard `slog` fields mapped to Cloud Logging's format.
*   `WithCloudTraceContext` HTTP middleware for `X-Cloud-Trace-Context` header integration, correlating logs with traces.
*   Custom log levels (e.g., `LevelNotice`, `LevelCritical`) aligned with Cloud Logging severities.

## Getting Started

1.  **Installation:**
    ```bash
    go get github.com/duizendstra/dui-go
    ```

2.  **Import & Use:**
    Import the desired packages into your Go application:
    ```go
    import (
        "context"
        "fmt"
        "log"
        "os"
        "time"

        "github.com/duizendstra/dui-go/env"
        "github.com/duizendstra/dui-go/authentication"
        "github.com/duizendstra/dui-go/cache"
        "github.com/duizendstra/dui-go/logging/cloudlogging" // For server/API logging
        "log/slog"
    )
    ```


## Basic Usage Example

This example demonstrates loading application configuration using the `env` package and setting up a basic token manager from the `authentication` package.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/duizendstra/dui-go/env"
	"github.com/duizendstra/dui-go/authentication"
	"github.com/duizendstra/dui-go/cache"
	// For a server application, you might use the cloudlogging handler:
	// "github.com/duizendstra/dui-go/logging/cloudlogging"
	// stdslog "log/slog" // Alias if using both log and slog
)

// AppConfig defines application configuration loaded from environment variables.
type AppConfig struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"MyApplication"`
	APIKey      string `env:"API_KEY" envRequired:"true"`
	APITimeout  int    `env:"API_TIMEOUT_SECONDS" envDefault:"10"`
}

// MockTokenFetcher simulates fetching a token from an external service.
func MockTokenFetcher(apiKey string) authentication.TokenFetcher {
	return func() (string, time.Time, error) {
		// In a real scenario, use apiKey to fetch a real token.
		log.Printf("Fetching new token using API key: ...%s\n", apiKey[len(apiKey)-4:])
		return fmt.Sprintf("mock-token-for-%s", apiKey), time.Now().Add(1 * time.Hour), nil
	}
}

func main() {
	// --- 1. Load Configuration ---
	// Set up example environment variables
	os.Setenv("API_KEY", "secret-key-12345")
	os.Setenv("API_TIMEOUT_SECONDS", "5")
	defer os.Unsetenv("API_KEY")
	defer os.Unsetenv("API_TIMEOUT_SECONDS")

	var cfg AppConfig
	if err := env.Process(&cfg); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Service '%s' configured with API timeout %d seconds.\n", cfg.ServiceName, cfg.APITimeout)

	// --- 2. Setup Authentication Token Manager ---
	// Create a cache (e.g., in-memory for this example)
	inMemCache := cache.NewInMemoryCache()

	// Create a TokenManager
	tokenManager := authentication.NewTokenManager(inMemCache)

	// Register a fetcher for our service's token
	// The fetcher might use configuration like cfg.APIKey
	tokenManager.RegisterFetcher("myExternalServiceToken", MockTokenFetcher(cfg.APIKey))

	// --- 3. Use the Token Manager ---
	ctx := context.Background() // Or your request context
	serviceToken, err := tokenManager.GetToken(ctx, "myExternalServiceToken") // Pass context to GetToken
	if err != nil {
		log.Fatalf("Failed to get service token: %v", err)
	}
	fmt.Printf("Successfully retrieved service token: %s\n", serviceToken)

	// Subsequent calls might use the cached token
	serviceToken2, _ := tokenManager.GetToken(ctx, "myExternalServiceToken")
	fmt.Printf("Retrieved service token again: %s\n", serviceToken2)

	// Example Output:
	// Service 'MyApplication' configured with API timeout 5 seconds.
	// Fetching new token using API key: ...345
	// Successfully retrieved service token: mock-token-for-secret-key-12345
	// Retrieved service token again: mock-token-for-secret-key-12345
}

```
*(Note: `TokenManager.GetToken` was updated in a previous step to accept `context.Context`. The example above reflects that. Ensure your `TokenManager.GetToken` signature matches.)*

For more detailed examples, please refer to the `_test.go` files within each package, particularly the `Example_...` functions.

## Building and Testing

To build the `dui-go` library locally (e.g., for contributing):
```bash
go build ./...
```

To run all tests for the library:
```bash
go test ./...
```

To run tests with the Go race detector (highly recommended):
```bash
go test -race ./...
```

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.MD](CONTRIBUTING.MD) for guidelines on how to submit pull requests, report issues, and other contribution details.

## Roadmap

To see what features are planned for future releases, please check the [ROADMAP.MD](ROADMAP.MD).

## Code of Conduct

This project and its community are expected to adhere to a standard of professional and respectful conduct. Please be kind, considerate, and welcoming in all interactions. Harassment or exclusionary behavior will not be tolerated.

## License

`dui-go` is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```