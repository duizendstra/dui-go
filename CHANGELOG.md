# Changelog

All notable changes to the **dui-go** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
*(For next version after 0.1.0)*

### Changed
*(For next version after 0.1.0)*

### Fixed
*(For next version after 0.1.0)*

---

## [0.1.0] - 2025-05-10

### Added
*   **Comprehensive AI Development Guidelines:**
    *   Core `.idx/airules.md` (v3.1) established for general Go AI assistance in Firebase Studio, instructing AI to prioritize project-specific rules.
    *   Dedicated `docs/AI_PROJECT_SPECIFICS.md` (v1.0) created with detailed AI guidelines, context, architectural patterns, and specific rules tailored for `dui-go` library development.
*   **Enhanced Contributor Documentation (`CONTRIBUTING.md`):**
    *   Added detailed "Documentation Standards and Workflow" section, including filename conventions and a release finalization checklist.
    *   Incorporated recommendations for using the `ContextVibes CLI` for development workflows.
*   **Firebase Studio Integration:**
    *   Added an "Explore in Firebase Studio" button and explanatory section to `README.md` for easier access to the `dui-go` codebase within the IDE.
    *   Updated `dui-go/.idx/dev.nix` to include `pkgs.gh` and an `onCreate` script to install a pinned version of `ContextVibes CLI` (v0.1.1 recommended, adjust if different), plus an `onStart` script to check for CLI presence.

### Changed
*   **Documentation Overhaul & Refinement:**
    *   Revised and enhanced all package-level `doc.go` files for improved clarity, purpose statements, and illustrative usage snippets (targeting an "Excellent" standard).
    *   Refined root `README.md` for better structure, added a "Why `dui-go`?" section, improved package linking (to subdirectories/godoc), and updated the main usage example.
    *   Moved `ROADMAP.md` to `docs/ROADMAP.md` and ensured its content reflects current plans.
    *   Updated `.gitignore` with more robust include-based logic, ensuring correct tracking of documentation and ignoring of CLI/OS artifacts (including `STRATEGIC_KICKOFF_PROTOCOL_FOR_AI.md` and user-specific `.contextvibes.yaml`).
*   **Minor `doc.go` updates:** Ensured all examples are illustrative and clear. (Though no major code changes to the library itself were made in this pass, the documentation represents a significant update to how the library is understood and developed with AI assistance).

---

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
[Unreleased]: https://github.com/duizendstra/dui-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/duizendstra/dui-go/compare/v0.0.2...v0.1.0
[0.0.2]: https://github.com/duizendstra/dui-go/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/duizendstra/dui-go/releases/tag/v0.0.1