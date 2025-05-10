# dui-go Library - Roadmap

This document outlines the planned features and future direction for the `dui-go` library. Priorities and features may change based on feedback and development progress. This roadmap aims to guide development efforts.

## Near Term (Next 1-2 Minor Releases)

*   [ ] **Implement single-flight for `authentication.TokenManager`:** *(High Priority)*
    *   Integrate `golang.org/x/sync/singleflight` into `TokenManager.GetToken` to prevent redundant concurrent fetches for the same token key (thundering herd problem).
    *   Goal: Improve performance and reduce load on underlying token sources under concurrent access patterns.
*   [ ] **Comprehensive Unit Testing for `authentication` Package:**
    *   Ensure `TokenManager` tests are fully deterministic by using mockable time for all expiry checks.
    *   Add specific tests for concurrent access to `TokenManager.GetToken` to verify single-flight behavior.
    *   Goal: Increase confidence in the concurrency and caching logic.
*   [ ] **Interface Review and Refinement for `cache` Package:**
    *   Evaluate the current `cache.Cache` interface.
    *   Consider if methods like `SetAll`, `GetAll`, `Flush` are universally applicable or if a more minimal core interface with optional extensions would be better for broader backend compatibility.
    *   Goal: Ensure the cache interface is both powerful and flexible for various implementations.
*   [ ] **Error Handling Consistency:**
    *   Review all packages to ensure consistent use of `fmt.Errorf` with `%w` for error wrapping, or the library's `errors` package where appropriate.
    *   Standardize error types returned by public APIs where beneficial.
    *   Goal: Improve debuggability and programmatic error handling for library users.

## Medium Term

*   [ ] **Expand `cache` Package Implementations:**
    *   Consider adding a Redis-backed implementation of `cache.Cache`.
    *   Investigate options for a distributed cache implementation.
*   [ ] **More Sophisticated `env` Package Features:**
    *   Explore adding support for `.env` file loading.
    *   Consider more advanced validation rules or custom validator registration.
*   [ ] **Context Propagation Review:**
    *   Ensure `context.Context` is consistently the first argument for all relevant functions and methods that might perform I/O or be long-running.
    *   Verify cancellation and deadline propagation.
*   [ ] **Benchmarking:**
    *   Add benchmarks for performance-sensitive components like `cache` implementations and `authentication.TokenManager`.
    *   Goal: Identify performance bottlenecks and validate optimizations.

## Long Term / Vision

*   [ ] **Observability Enhancements:**
    *   Explore adding optional OpenTelemetry (OTel) instrumentation hooks within key packages (e.g., `TokenManager` fetches, cache operations, Firestore interactions).
    *   Goal: Allow library users to easily integrate `dui-go` components into their existing observability stacks.
*   [ ] **Additional Utility Packages:**
    *   Identify other common cross-cutting concerns in Go/GCP development that could benefit from a utility package (e.g., retry mechanisms, common GCP client initializations).
*   [ ] **Configuration Management Package:**
    *   Consider a more comprehensive configuration management package that could load from environment variables, files (YAML, JSON, TOML), and potentially integrate with GCP Secret Manager.

---

This roadmap is a living document. Feedback and suggestions are highly welcome!