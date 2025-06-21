# Changelog

All notable changes to the **dui-go** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
*(For next version after 0.3.0)*

### Changed
*(For next version after 0.3.0)*

### Fixed
*(For next version after 0.3.0)*

---

## [0.3.0] - 2025-05-21

This release introduces a new Google Cloud Storage client and includes a major investment in overall project quality, testing, and documentation.

### Added
*   **`gcs` Package:** Introduced a new package for simplified interaction with Google Cloud Storage, featuring a robust, streaming `Upload` method.
*   **Backlog Documentation:** Created formal backlog items in `docs/backlog/` to track future technical debt and improvement tasks for testing and refactoring.

### Changed
*   **Standardized Client APIs:** The `gcs` client was updated to use a `Config` struct for initialization, making its API consistent with the `secretmanager` client for a better developer experience.
*   **Refined Logging Strategy:** Tuned logging in both `gcs` and `secretmanager` clients. Successful operations now log at the `DEBUG` level (instead of `INFO`) to reduce noise, and all logging calls now use context-aware variants (e.g., `slog.DebugContext`) to include trace information.
*   **Improved Test Structure:** Refactored tests in `authentication/tokenmanager_test.go` to use `t.Run()` for named sub-tests, making test output clearer and removing explanatory comments.
*   **Overhauled `README.md`:** The main project README was rewritten to be more comprehensive, with a clear project philosophy, a complete package overview, and improved examples.
*   **Renamed AI Guide:** The project-specific AI guide was renamed from `docs/AI_PROJECT_SPECIFICS.md` to `docs/ai_context.md` for clarity and convention. All internal references were updated.

### Fixed
*   **Corrected `.gitignore`:** The `.gitignore` file was updated to correctly track the shared `/.vscode/settings.json` file while ignoring other user-specific workspace files.
*   **Cleaned Test Comments:** Removed historical and placeholder comments from several test files (e.g., `firestorekv_test.go`) to ensure they accurately reflect the current state of the codebase.

## [0.2.0] - 2025-05-21

### Added
*   **`secretmanager` Package:** Added a new client for securely fetching secrets from Google Cloud Secret Manager. The client is initialized via a `Config` struct and includes optional `slog.Logger` support for structured logging.

---

## [0.1.0] - 2025-05-10

### Added
*   **Comprehensive AI Development Guidelines:**
    *   Core `.idx/airules.md` (v3.1) established for general Go AI assistance in Firebase Studio.
    *   Dedicated `docs/AI_PROJECT_SPECIFICS.md` (v1.0) created with detailed AI guidelines tailored for `dui-go` library development.
*   **Enhanced Contributor Documentation (`CONTRIBUTING.md`):**
    *   Added detailed "Documentation Standards and Workflow" section.
*   **Firebase Studio Integration:**
    *   Added an "Explore in Firebase Studio" button and section to `README.md`.

### Changed
*   **Documentation Overhaul & Refinement:**
    *   Revised and enhanced all package-level `doc.go` files for improved clarity.
    *   Refined root `README.md` for better structure and updated usage examples.
    *   Moved `ROADMAP.md` to `docs/ROADMAP.md`.

---

## [0.0.2] - 2025-05-07

### Changed
*   **BREAKING CHANGE:** Refactored the `env` package for type-safe environment variable loading into structs using tags.

---

## [0.0.1] - 2025-05-07

### Added
*   Initial public release of the `dui-go` utility library.
*   Core packages: `authentication`, `cache`, `env`, `errors`, `firestore`, `logging/cloudlogging`, `store`.
*   Internal testing utilities in `internal/testutil`.
*   Project foundation files (`go.mod`, `README.md`, `LICENSE`, etc.).

---

<!-- Link Definitions -->
[Unreleased]: https://github.com/duizendstra/dui-go/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/duizendstra/dui-go/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/duizendstra/dui-go/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/duizendstra/dui-go/compare/v0.0.2...v0.1.0
[0.0.2]: https://github.com/duizendstra/dui-go/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/duizendstra/dui-go/releases/tag/v0.0.1