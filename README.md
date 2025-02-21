# go-dui

`go-dui` is a foundational library of utilities designed to streamline the development of Go-based services, particularly within the context of the EasyFlor project and Google Cloud Platform. It provides essential components for common tasks such as authentication, caching, environment variable management, structured error handling, key-value data storage, and logging.

## Packages

The `go-dui` library is organized into several packages, each addressing a specific area of functionality:

### Authentication (`go-dui/authentication`)

This package provides a robust mechanism for managing authentication tokens. It defines interfaces for `Token`, `TokenManager`, and `TokenFetcher`, enabling flexible token retrieval, caching, and refresh strategies.

**Key Features:**

*   **`Token` struct:** Represents a token with its value and expiration time.
*   **`TokenManager` interface:** Defines methods for retrieving and managing tokens.
*   **`TokenFetcher` interface:** Abstracts the process of fetching tokens from an external source.
*   **Thread-safe `TokenManager` implementation:** Ensures safe concurrent access to tokens.

### Cache (`go-dui/cache`)

This package offers a generic `Cache` interface for in-memory and external caching. It supports customizable expiration policies and provides a basic in-memory cache implementation.

**Key Features:**

*   **`Cache` interface:** Defines basic `Get` and `Set` operations for caching.
*   **`Options` and `Option`:** Allow configuration of cache behavior, including expiration.
*   **`inMemoryCache`:** A thread-safe in-memory cache implementation.

### Env (`go-dui/env`)

This package simplifies environment variable handling by providing type-safe accessors and default value support. It allows defining namespaces for organizing environment variables.

**Key Features:**

*   **`Namespace` type:** Represents a namespace for environment variables.
*   **Type-safe accessors:** `Get()`, `Int()`, `Bool()`, and their variants for different data types.
*   **Default value support:** `GetDefault()`, `IntDefault()`, `BoolDefault()` for specifying default values.
*   **`Must` variants:** `MustGet()`, `MustInt()`, `MustBool()` for panicking on missing or invalid values.

### Errors (`go-dui/errors`)

This package enhances Go's error handling by introducing error codes and associating them with HTTP status codes. It provides a convenient way to categorize errors and handle them appropriately in API responses.

**Key Features:**

*   **`Error` interface:** Extends the standard `error` interface with a `Code()` method.
*   **`Errorf()`:** Creates formatted errors with codes.
*   **`Is()`:** Checks if an error is of a specific code.
*   **`StatusCode()`:** Returns the HTTP status code associated with an error code.

### Firestore (`go-dui/firestore`)

This package provides a key-value store abstraction over Google Cloud Firestore. It enables storing and retrieving arbitrary data structures by serializing them as JSON. It also implements the `cache.Cache` interface, allowing Firestore to be used as a cache backend.

**Key Features:**

*   **`FirestoreKV`:** Implements `cache.Cache` for using Firestore as a cache.
*   **JSON serialization:** Stores and retrieves values as JSON.
*   **Retry mechanism:** Enhances robustness against transient errors.

### Store (`go-dui/store`)

This package defines a generic key-value `Store` interface and provides a `FirestoreStore` implementation that uses Firestore as the backend.

**Key Features:**

*   **`Store` interface:** Defines basic `Get`, `Set`, and `Delete` operations.
*   **`FirestoreStore`:** Implements `Store` using Firestore.
*   **Generic `Document` type:** Supports flexible data representation.

### Logging/Cloudlogging (`go-dui/logging/cloudlogging`)

This package provides a logging handler specifically designed for Google Cloud Logging. It formats log entries according to Cloud Logging's structured logging format and maps standard log levels to Cloud Logging severities.

**Key Features:**

*   **`NewHandler()`:** Creates a logging handler for Cloud Logging.
*   **Structured logging:** Formats log entries for optimal Cloud Logging integration.
*   **Severity mapping:** Maps log levels to Cloud Logging severities.
*   **`NewMiddleware()`:** Creates a middleware for adding request details to log context.

## Getting Started

1. **Installation:**

    ```bash
    go get github.com/path/to/go-dui # Replace with the actual path to the go-dui repository
    ```

2. **Usage:**

    Import the desired packages into your Go code:

    ```go
    import (
        "github.com/path/to/go-dui/authentication"
        "github.com/path/to/go-dui/cache"
        "github.com/path/to/go-dui/env"
        "github.com/path/to/go-dui/errors"
        "github.com/path/to/go-dui/firestore"
        "github.com/path/to/go-dui/store"
        "github.com/path/to/go-dui/logging/cloudlogging"
    )
    ```

    Refer to the individual package documentation and examples for detailed usage instructions.

## Building and Testing

To build the `go-dui` library, run:

```bash
task build
