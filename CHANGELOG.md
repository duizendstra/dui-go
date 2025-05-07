# Changelog

All notable changes to the **dui-go** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

## [0.0.2] - 2025-05-07

### Changed
*   **BREAKING CHANGE:** Refactored the `env` package for type-safe environment variable loading into structs.
    *   Removed the previous API based on `env.EnvVar` slices and the `env.EnvLoader` which returned `map[string]string`.
    *   Introduced the `env.Process(spec any)` function which populates the fields of a given struct pointer based on environment variables.
    *   Configuration is now managed via struct tags:
        *   `env:"VAR_NAME"`: Specifies the environment variable name (defaults to uppercase field name).
        *   `envDefault:"value"`: Provides a default string value if the environment variable is not set or empty.
        *   `envRequired:"true"`: Marks the environment variable as mandatory.
    *   Supported field types for parsing: `string`, `int`, `int64`, `bool`.

## [0.0.1] - 2025-05-07

### Added
*   Initial public release of the `dui-go` utility library.
*   Core packages:
    *   `authentication`: Token management.
    *   `cache`: In-memory caching interface and implementation.
    *   `env`: Environment variable loading and validation (map-based).
    *   `errors`: Structured API error types.
    *   `firestore`: Firestore-backed key-value store.
    *   `logging/cloudlogging`: Google Cloud Logging `slog` handler and middleware.
    *   `store`: Generic key-value store interface with Firestore implementation.
*   **Internal Testing Utilities:**
    *   `internal/testutil`: Created mock implementations for `cache.Cache` and `firestore.KV` to support testing.
*   **Project Foundation:**
    *   Established Go module structure (`go.mod`, `go.sum`).
    *   Included essential documentation (`README.md`, `CONTRIBUTING.md`) and `LICENSE`.
    *   Added initial unit tests and example tests for core packages.
    *   Set up Project IDX configuration files (`.idx/dev.nix`, `.idx/airules.md`).

---

<!-- Link Definitions -->
[Unreleased]: https://github.com/duizendstra/dui-go/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/duizendstra/dui-go/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/duizendstra/dui-go/releases/tag/v0.0.1
