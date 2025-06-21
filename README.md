# dui-go

[![Go Reference](https://pkg.go.dev/badge/github.com/duizendstra/dui-go.svg)](https://pkg.go.dev/github.com/duizendstra/dui-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/duizendstra/dui-go)](https://goreportcard.com/report/github.com/duizendstra/dui-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`dui-go` is a foundational library of idiomatic Go utilities designed to accelerate the development of robust, maintainable, and testable services, with a focus on seamless Google Cloud Platform (GCP) integration.

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

## Project Philosophy: Why `dui-go`?

This library is built on a few core principles to make your development life easier:

*   **Reduce Boilerplate:** Provides well-tested, reusable components for common tasks like configuration, caching, and authentication, letting you focus on your application's core logic.
*   **Idiomatic Go & Best Practices:** Encourages idiomatic Go patterns, including clear API design, explicit error handling, and effective concurrency.
*   **Simplified GCP Integration:** Offers focused, easy-to-use clients for services like Cloud Storage, Secret Manager, and Firestore, including built-in trace context propagation for observability.
*   **Testability by Design:** Promotes a clean separation of concerns through interfaces (`cache.Cache`, `store.Store`), making your code inherently testable and allowing you to swap implementations with ease.

## Packages

The `dui-go` library is organized into several packages, each addressing a specific area of functionality.

| Package | Description |
| :--- | :--- |
| **[Authentication](./authentication/)** | A robust, thread-safe manager for the entire token lifecycle: retrieval, caching, and refresh logic. |
| **[Cache](./cache/)** | A flexible caching abstraction (`cache.Cache`) with a thread-safe in-memory implementation. |
| **[Env](./env/)** | Type-safe environment variable loading into Go structs using simple struct tags (`env`, `envDefault`, `envRequired`). |
| **[Errors](./errors/)** | Structured `APIError` types with codes and details, ideal for building consistent API error responses. |
| **[Firestore](./firestore/)** | A simplified key-value store abstraction (`firestore.KV`) built on top of Google Cloud Firestore. |
| **[GCS](./gcs/)** | A focused client for Google Cloud Storage, featuring a robust, streaming `Upload` method. |
| **[Logging/Cloudlogging](./logging/cloudlogging/)** | A `log/slog` handler for Google Cloud Logging that automatically formats logs and propagates trace context. |
| **[SecretManager](./secretmanager/)** | A secure client for fetching secrets from Google Cloud Secret Manager. |
| **[Store](./store/)** | A generic key-value `Store` interface that decouples your application from specific storage backends. |

## Getting Started

1.  **Installation:**
    ```bash
    go get github.com/duizendstra/dui-go
    ```

2.  **Basic Usage Example:**
    This example demonstrates loading configuration with `env` and managing tokens with `authentication`.

    ```go
    package main

    import (
    	"context"
    	"fmt"
    	"log"
    	"os"
    	"time"

    	"github.com/duizendstra/dui-go/authentication"
    	"github.com/duizendstra/dui-go/cache"
    	"github.com/duizendstra/dui-go/env"
    )

    // AppConfig defines our application's configuration.
    type AppConfig struct {
    	ServiceName string `env:"SERVICE_NAME" envDefault:"MyApplication"`
    	APIKey      string `env:"API_KEY" envRequired:"true"`
    	APITimeout  int    `env:"API_TIMEOUT_SECONDS" envDefault:"10"`
    }

    // MockTokenFetcher simulates fetching a token from a remote service.
    func MockTokenFetcher(apiKey string) authentication.TokenFetcher {
    	return func(ctx context.Context) (string, time.Time, error) {
    		// The TokenFetcher uses the API key loaded from config.
    		log.Printf("Fetching new token using API key: ...%s\n", apiKey[len(apiKey)-4:])
    		return fmt.Sprintf("mock-token-for-%s", apiKey), time.Now().Add(1 * time.Hour), nil
    	}
    }

    func main() {
    	// --- 1. Load Configuration from Environment ---
    	os.Setenv("API_KEY", "secret-key-12345")
    	defer os.Unsetenv("API_KEY")

    	var cfg AppConfig
    	if err := env.Process(&cfg); err != nil {
    		log.Fatalf("Error loading config: %v", err)
    	}
    	fmt.Printf("Service '%s' configured.\n", cfg.ServiceName)

    	// --- 2. Setup Dependencies & Token Manager ---
    	// The TokenManager depends on a cache. We provide an in-memory one.
    	inMemCache := cache.NewInMemoryCache()
    	tokenManager := authentication.NewTokenManager(inMemCache)

    	// Register the fetcher for the token we need.
    	tokenManager.RegisterFetcher("myExternalServiceToken", MockTokenFetcher(cfg.APIKey))

    	// --- 3. Use the Token Manager in Application Logic ---
    	ctx := context.Background()
    	// The manager handles caching and fetching transparently.
    	serviceToken, err := tokenManager.GetToken(ctx, "myExternalServiceToken")
    	if err != nil {
    		log.Fatalf("Failed to get service token: %v", err)
    	}
    	fmt.Printf("Successfully retrieved service token: %s\n", serviceToken)

    	// This second call will hit the cache; the fetcher log won't appear again.
    	serviceToken2, _ := tokenManager.GetToken(ctx, "myExternalServiceToken")
    	fmt.Printf("Retrieved service token again: %s\n", serviceToken2)
    }
    ```

## Building and Testing

This project maintains a high standard of quality and testing.

*   **Build the library:**
    ```bash
    go build ./...
    ```
*   **Run all unit tests:**
    ```bash
    go test ./...
    ```
*   **Run tests with the race detector (highly recommended):**
    ```bash
    go test -race ./...
    ```

## Contributing

Contributions are what make the open-source community such an amazing place. We welcome any contributions that improve this library.

Please read our **[CONTRIBUTING.md](CONTRIBUTING.md)** for details on our code of conduct, development process, and how to submit pull requests.

## Roadmap

To see what features are planned for future releases, please check the **[docs/roadmap.md](docs/roadmap.md)**.

## License

`dui-go` is distributed under the MIT License. See the **[LICENSE](LICENSE)** file for more information.