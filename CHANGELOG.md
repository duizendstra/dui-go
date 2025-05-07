# Changelog

All notable changes to the **dui-go** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.0.1] - 2025-05-07

### Added
- Initial public release of the `dui-go` utility library.
- Core packages:
    - `authentication`: Token management.
    - `cache`: In-memory caching interface and implementation.
    - `env`: Environment variable loading and validation.
    - `errors`: Structured API error types.
    - `firestore`: Firestore-backed key-value store.
    - `logging/cloudlogging`: Google Cloud Logging `slog` handler and middleware.
    - `store`: Generic key-value store interface with Firestore implementation.
*   **Internal Testing Utilities:**
    *   `internal/testutil`: Created mock implementations for `cache.Cache` and `firestore.KV` to support testing.
*   **Project Foundation:**
    *   Established Go module structure (`go.mod`, `go.sum`).
    *   Included essential documentation (`README.md`, `CONTRIBUTING.md`) and `LICENSE`.
    *   Added initial unit tests and example tests for core packages.
    *   Set up Project IDX configuration files (`.idx/dev.nix`, `.idx/airules.md`).

---

<!--
Link Definitions - Add the new one when tagging a new version.
The first entry (Unreleased) would typically look like this:
## [Unreleased]

[Unreleased]: https://github.com/duizendstra/dui-go/compare/v0.0.1...HEAD
-->
[0.0.1]: https://github.com/duizendstra/dui-go/releases/tag/v0.0.1