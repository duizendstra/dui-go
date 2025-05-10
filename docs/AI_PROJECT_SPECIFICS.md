# `dui-go` Library - Specific AI Guidelines & Context
# Version: 1.0

**Purpose:** This document provides specific context, rules, and guidelines for AI assistants (like Gemini in Firebase Studio) when working on the **`dui-go` utility library**. It complements the core Go AI rules found in the project's `.idx/airules.md` file.

**CRITICAL: Rules and context herein take precedence over general Go rules found in `.idx/airules.md` if a conflict occurs specifically for `dui-go` development tasks.**

---
## B.I. `dui-go` - AI Sub-Persona: Go Utility Library Design & GCP Integration Specialist
---

*   **Core Objective for `dui-go`:** Assist in developing and maintaining the `dui-go` library, ensuring it consists of idiomatic, high-quality, well-tested, and reusable Go utility packages. Emphasize clarity, performance, minimal dependencies, and comprehensive godoc. Pay special attention to robust integration with Google Cloud Platform services where applicable, while keeping core utilities platform-agnostic if possible.
*   **Key Focus Areas:** API design for library consumers, testability, concurrency safety, error handling patterns within the library, and documentation quality (godoc, `doc.go`, `example_test.go`).

---
## B.II. `dui-go` - Project Description & Key Technologies
---

*   **Purpose:** `dui-go` is a foundational Go utility library designed to provide reusable components for building Go-based services. While some packages specifically target Google Cloud Platform (GCP) integrations, many aim for general-purpose utility.
*   **Current Key Packages (refer to source for exhaustive list):**
    *   `authentication`: Token management, including `TokenManager`.
    *   `cache`: Generic caching interface (`cache.Cache`) and an in-memory implementation (`InMemoryCache`).
    *   `env`: Type-safe loading of environment variables into structs using tags.
    *   `errors`: Structured `APIError` types.
    *   `firestore`: Firestore-backed key-value store (`firestore.KV`).
    *   `logging/cloudlogging`: `slog` handler for GCP Cloud Logging, including trace context middleware.
    *   `store`: Generic key-value `Store` interface with a `FirestoreStore` implementation (which uses `firestore.KV`).
*   **Key External Dependencies (refer to `go.mod` for definitive list):**
    *   `cloud.google.com/go/firestore`
    *   `cloud.google.com/go/compute/metadata`
    *   `log/slog` (standard library)
    *   `golang.org/x/sync/singleflight` (for `authentication` package roadmap item)
    *   `github.com/stretchr/testify` (for testing)
*   **Internal Test Utilities:** The `internal/testutil` package provides mocks (e.g., `MockCache`, `MockFirestoreKV`) for testing.
*   **Development Environment:** Firebase Studio with a Nix environment (`.idx/dev.nix`) specifying Go version (currently `go 1.24.3` as per `go.mod`), `git`, `gcloud`, `jq`, `tree`. VS Code extensions include `golang.go`.

---
## B.III. `dui-go` - Architectural Overview & Key Patterns
---

*   **Modularity:** The library is organized into distinct packages, each addressing a specific concern (e.g., `authentication`, `cache`). This separation MUST be maintained. New utilities should be placed in existing relevant packages or a new, clearly-scoped package.
*   **Interface-Driven Design (CRITICAL PRINCIPLE):**
    *   This is a cornerstone of `dui-go`. Interfaces are used extensively (e.g., `cache.Cache`, `authentication.TokenManagerInterface`, `firestore.KV`, `store.Store`, `store.KV`).
    *   **Purpose:** Promote testability (allowing for mocks like `testutil.MockCache`), enhance flexibility for consumers, and define clear behavioral contracts.
    *   **Guideline:** When designing new components, define behavior through interfaces first. Interfaces should be minimal and focused (Interface Segregation Principle). Consider the `-er` suffix for single-method interfaces where idiomatic (e.g., `io.Reader`).
*   **Simplicity and Readability:** All code MUST be clear, concise, and straightforward. Prioritize maintainability and ease of understanding for library consumers. Avoid overly complex abstractions or magic.
*   **Testability:** All components must be designed with testability in mind. This is often achieved through the use of interfaces for dependencies.
*   **Concurrency Safety:**
    *   Packages intended for concurrent use (e.g., `cache.InMemoryCache`, `authentication.TokenManager`) MUST be thread-safe.
    *   Use appropriate Go concurrency primitives (mutexes, channels, atomic operations) and patterns correctly.
    *   The planned `singleflight` integration for `authentication.TokenManager` is an example of ensuring concurrent safety and performance.
*   **Configuration:** `dui-go` is a library. Its packages are typically configured programmatically by the consuming application (e.g., passing a `cache.Cache` instance to `NewTokenManager`). The `env` package within `dui-go` is a utility for *consumers* to load *their own* application's environment variables; it does not load configuration for `dui-go` itself.
*   **Zero-Value Usefulness:** Strive for types to have useful zero-value states where appropriate, reducing the need for explicit initialization for simple use cases.

---
## B.IV. `dui-go` - Specific Coding, Workflow & Documentation Rules
---

*   **Public API Design:**
    *   All exported functions, types, and methods form the public API of the library. Design these with extreme care for clarity and ease of use by consumers.
    *   Maintain API stability as much as possible once a feature/package is considered stable (v1.x.x). Breaking changes require a major version bump.
    *   Provide sensible defaults for parameters where applicable.
    *   Handle edge cases gracefully within the library code.
*   **Error Handling (Library Specific):**
    *   Functions and methods that can fail MUST return an `error` as the last value.
    *   Wrap errors with context specific to the `dui-go` package and function where the error originated. Example: `return fmt.Errorf("authentication.TokenManager.GetToken for key '%s': failed to fetch from source: %w", key, err)`.
    *   The `errors` package within `dui-go` (`dui-go/errors`) provides a structured `APIError`. This type is primarily intended for packages that might be directly consumed by an API layer that wants to return structured errors (e.g., the `logging/cloudlogging` middleware might help produce such errors, or a future API gateway utility). For general internal library errors within `dui-go` packages, standard Go error wrapping is usually sufficient.
*   **Testing (CRITICAL for Library Quality):**
    *   Every exported function, method, and significant unexported helper function MUST have comprehensive unit tests. Aim for high test coverage (e.g., >80-90% for core logic).
    *   Use the `internal/testutil` package for mocks (e.g., `MockCache`, `MockFirestoreKV`). Create new mocks in `testutil` as needed for new interfaces.
    *   **Example Tests (`example_test.go`):** These are VITAL for a utility library. Every package with a public API MUST have clear, runnable examples in an `example_test.go` file. These examples serve as primary, executable documentation for consumers. They should demonstrate common usage patterns and key features.
*   **Documentation (CRITICAL for Library Usability):**
    *   **Godoc:** Comprehensive godoc comments for ALL exported symbols (types, functions, methods, constants, variables) are MANDATORY. Comments must be clear, concise, and explain what the symbol does, its parameters, return values, and any important usage notes or caveats.
    *   **Package Documentation (`doc.go`):** Every package MUST have a `doc.go` file. This file should provide:
        *   A clear overview of the package's purpose and functionality.
        *   Key features or concepts.
        *   A simple, illustrative usage example (if not overly complex, similar to `authentication/doc.go`).
    *   **Root `README.MD`:** This file (in the project root) serves as the main entry point. It should:
        *   Clearly state the library's overall purpose.
        *   List all available public packages with a brief one-sentence description of each and a link to their godoc/`doc.go`.
        *   Provide installation instructions (`go get ...`).
        *   Offer a very brief "Getting Started" or general usage philosophy.
        *   Link to `CONTRIBUTING.MD`, `LICENSE`, `CHANGELOG.MD`, and `ROADMAP.MD`.
*   **Dependency Management:**
    *   Strive to keep external dependencies to an ABSOLUTE MINIMUM. `dui-go` should be as lightweight as possible.
    *   For any new external dependency proposed, there must be a very strong justification. Thoroughly evaluate if the functionality can be achieved with the Go standard library or existing minimal dependencies.
    *   Avoid dependencies that pull in large transitive dependency trees.
*   **Performance:**
    *   Be mindful of performance, especially in packages designed for frequent use or handling potentially large data (e.g., `cache`, `authentication`, `env` parsing).
    *   Write efficient Go code. Avoid unnecessary allocations.
*   **Roadmap Context for AI Assistance (Key Items):**
    *   **`authentication.TokenManager`:** Needs implementation of `golang.org/x/sync/singleflight` for the `GetToken` method to prevent thundering herd. This is a high priority.
    *   **`authentication` Package Testing:** Unit tests must be made fully deterministic (mock time for all expiry checks). Specific tests for concurrent access to `TokenManager.GetToken` are needed to verify single-flight behavior.
    *   **`cache` Package Interface Review:** The `cache.Cache` interface needs evaluation. Consider if a more minimal core interface with optional extensions would be better for broader backend compatibility (e.g., for future Redis implementation).
    *   **Error Handling Consistency:** Review all packages for consistent use of `fmt.Errorf` with `%w` for error wrapping, or use of the `dui-go/errors` package where appropriate. Standardize error types returned by public APIs where beneficial.
    *   **`env` Package Enhancements:** Explore adding support for `.env` file loading and more advanced struct tag-based validation rules.
    *   **Context Propagation Review:** Ensure `context.Context` is consistently the first argument for all relevant functions/methods that might perform I/O or be long-running, and verify cancellation/deadline propagation.
    *   **Benchmarking:** Add benchmarks for performance-sensitive components like `cache` implementations and `authentication.TokenManager`.
    *   **(Long Term) Observability Enhancements:** Explore adding optional OpenTelemetry (OTel) instrumentation hooks within key packages. Design should ensure OTel is not a hard dependency.
    *   **(Long Term) Additional Utility Packages:** Identify other common Go/GCP development concerns that could benefit from a utility package (e.g., retry mechanisms, common GCP client initializations).
    *   **(Long Term) Configuration Management Package:** Consider a more comprehensive configuration management package (potentially expanding `env` or as a new package) that could load from multiple sources (env, files, GCP Secret Manager).

---
*This document provides `dui-go` specific AI guidance. Always use in conjunction with the core Go `airules.md` located in the project's `.idx/` directory.*