# Contributing to dui-go

Thank you for considering contributing to `dui-go`! We welcome improvements, bug fixes, and new features that help create a robust and useful Go utility library. Your contributions are valuable in making `dui-go` a high-quality resource for Go developers.

## Code of Conduct

All contributors and participants are expected to act professionally and respectfully. Please be kind, considerate, and welcoming in all interactions. Harassment or exclusionary behavior will not be tolerated. We aim to foster an open and inclusive environment.

## Getting Started

1.  **Prerequisites:**
    *   Ensure you have **Go** installed and configured correctly (refer to the `go.mod` file for the specific version, currently `1.24` or later recommended).
    *   **Git** installed and configured.
    *   Familiarity with Go development best practices.

2.  **Fork & Clone:**
    *   Fork the repository on GitHub: `https://github.com/duizendstra/dui-go` (or your correct upstream URL).
    *   Clone your fork locally:
        ```bash
        # Replace YOUR_USERNAME with your actual GitHub username
        git clone https://github.com/YOUR_USERNAME/dui-go.git dui-go-fork
        cd dui-go-fork
        ```

3.  **Set Upstream Remote:**
    It's good practice to configure the original repository as an upstream remote:
    ```bash
    git remote add upstream https://github.com/duizendstra/dui-go.git
    ```

4.  **Verify Setup:**
    Ensure you can build the library components and run tests:
    ```bash
    go build ./...
    go test ./...
    ```

## Recommended Tooling for Workflow Management

For a streamlined development workflow (managing branches, commits, context for AI, etc.), we recommend using the **ContextVibes CLI**.
*   **Installation & Usage:** See [github.com/contextvibes/cli](https://github.com/contextvibes/cli)
*   Using `contextvibes kickoff`, `contextvibes commit`, `contextvibes sync`, etc., can help maintain consistency with project conventions.

## Making Changes

1.  **Branching Strategy:**
    *   Always create a new branch from the latest `main` branch for your changes.
    *   Keep your `main` branch in sync with `upstream/main`.
        ```bash
        # Using standard Git:
        git checkout main
        git fetch upstream
        git rebase upstream/main
        git checkout -b feature/your-descriptive-feature-name
        # or for fixes:
        # git checkout -b fix/issue-description-or-number

        # Alternatively, if using ContextVibes CLI (after strategic kickoff is marked complete for this project):
        # contextvibes kickoff --branch feature/your-descriptive-feature-name
        ```
    *   Branch names should be descriptive, using lowercase letters and hyphens (e.g., `feature/add-redis-cache-adapter`, `fix/token-manager-expiry-bug`). If using ContextVibes CLI, it may enforce patterns defined in a local `.contextvibes.yaml`.

2.  **Implementation:**
    *   Make your code changes, keeping them focused on a single feature, bug fix, or improvement per branch/Pull Request.
    *   Write clean, idiomatic Go code.

3.  **Follow Style and Quality Standards:**
    *   **Formatting:** All Go code MUST be formatted with `gofmt`. Run `go fmt ./...` (or `contextvibes format` if using the CLI) before committing.
    *   **Linting:**
        *   Run `go vet ./...` to catch common issues.
        *   If `golangci-lint` is configured for the project (check for a `.golangci.yml` file), ensure your code passes its checks: `golangci-lint run ./...`.
        *   (You can also use `contextvibes quality` which may incorporate these checks).
    *   **AI Assistance:** If using an AI assistant (like Gemini in Firebase Studio) for code generation or suggestions:
        *   Ensure the AI is guided by the rules in `.idx/airules.md` (the core Go AI rules).
        *   Ensure the AI is also aware of and prioritizes the project-specific guidelines in `docs/AI_PROJECT_SPECIFICS.MD`.
        *   Critically review all AI-generated code for correctness, performance, security, and adherence to project standards before committing.
        *   (The `contextvibes describe` command can help prepare context for AI assistants).

4.  **Testing (CRITICAL):**
    *   **Unit Tests:** All new functions and methods, and any modified ones, MUST have corresponding unit tests in `_test.go` files. Contributions that increase overall test coverage are highly encouraged.
    *   **Example Tests:** For new public APIs or significant changes to existing ones in this library, `example_test.go` files demonstrating usage are VITAL.
    *   Run tests frequently:
        ```bash
        go test ./...
        # ALWAYS run tests with the race detector before submitting a PR:
        go test -race ./...
        # (Or use `contextvibes test -race ./...`)
        ```
    *   Ensure all tests pass.

5.  **Commit Messages:**
    *   Write clear, concise, and descriptive commit messages.
    *   We recommend following the [Conventional Commits](https://www.conventionalcommits.org/) specification. This helps in generating automated changelogs and provides a clear history.
        *   Examples:
            *   `feat(cache): Add Redis adapter for Cache interface`
            *   `fix(auth): Correct token expiry logic in TokenManager`
            *   `docs(readme): Update installation instructions`
            *   `refactor(store): Simplify Get method in FirestoreStore`
            *   `test(env): Add more cases for struct tag parsing`
    *   Commit messages should be written in the imperative mood (e.g., "Add feature" not "Added feature" or "Adds feature").
    ```bash
    git add . # Stage your changes
    git commit -m "feat(cache): Implement SetAll method for InMemoryCache"
    # (Or use `contextvibes commit -m "..."` which may apply configured validation)
    ```

## Documentation Standards and Workflow

Maintaining clear, consistent, and up-to-date documentation is crucial for `dui-go`.

### Filename Naming Conventions

*   **Primary Human-Readable Project Docs:** Files in the project root (e.g., `README.MD`, `LICENSE`, `CONTRIBUTING.MD`, `CHANGELOG.MD`, `ROADMAP.MD`) and key guides within the `docs/` folder (e.g., `docs/AI_PROJECT_SPECIFICS.MD`) should use `UPPERCASE_WITH_UNDERSCORES.MD` (e.g., `AI_PROJECT_SPECIFICS.MD`).
*   **Configuration-like Files:** Files like `.idx/airules.md` or `.golangci.yml` should use their standard lowercase naming.
*   **Go Source Files:** `lowercase_with_underscores.go` or `camelCase.go` as per standard Go conventions (typically `lowercase.go` for package files and `lowercase_with_underscores_test.go` for tests).

### Content & Style

*   Write documentation clearly and concisely.
*   Use Markdown effectively for formatting.
*   **Godoc:** All exported Go symbols (packages, types, functions, methods, constants) MUST have comprehensive Godoc comments.
*   **AI Guidelines:** When using AI assistance for development, refer to `.idx/airules.md` (core Go rules) and `docs/AI_PROJECT_SPECIFICS.MD` (`dui-go` specific rules) to guide the AI.

### Documentation Update Workflow

#### Initial Setup / New Major Components
When establishing documentation for a new major component or during initial project setup, we generally follow this order to ensure dependencies are met:
1.  `LICENSE`
2.  Core `.idx/airules.md`
3.  Project-specific `docs/AI_PROJECT_SPECIFICS.MD`
4.  `ROADMAP.MD`
5.  `CONTRIBUTING.MD` (this file!)
6.  Code-level docs (`doc.go`, Godoc comments, `example_test.go` files)
7.  `README.MD`
8.  `CHANGELOG.MD` (prepare structure with an "Unreleased" section)

#### Ongoing Updates
*   **Code-Level Documentation (`doc.go`, Godoc, `example_test.go`):** MUST be updated concurrently with any changes to public APIs, types, or functions. These are part of the code itself.
*   **`README.MD`**: Update whenever user-facing aspects of the library change significantly (e.g., addition of new core packages, major API changes that affect usage, changes to installation instructions).
*   **`docs/AI_PROJECT_SPECIFICS.MD`**: Update if `dui-go`'s specific architecture, key patterns, or focus areas for AI assistance change.
*   **`.idx/airules.md`**: This core file is expected to be more stable. Updates would typically involve refining general Go AI best practices.

#### Finalizing a Release (Order of Document Updates)

When preparing for a new version release, the following documents are typically updated in this order, after all code changes for the version are complete, tested, and merged into the release branch (or `main` if releasing from there):

| Order | Document                                  | Action & Purpose                                                                                                                               |
|-------|-------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| 1     | Code-level Docs (Godoc, `example_test.go`) | Ensure all new/modified code is fully documented and examples are accurate. This should ideally be done *during* development.                 |
| 2     | `README.MD`                               | Update with any new packages, significant features, or changes to installation/usage relevant to the upcoming release.                         |
| 3     | `ROADMAP.MD`                              | Mark completed items from the roadmap that are part of this release. Review and update near-term/medium-term goals if priorities have shifted. |
| 4     | `CHANGELOG.MD`                            | **This is a critical step before tagging.** Add a new version heading. Summarize all notable changes (Added, Changed, Fixed, Removed, Security) for this specific version. Be thorough and clear. |
| 5     | `CONTRIBUTING.MD` (If needed)             | Review if any development processes, standards, or tooling mentioned here have changed with this release cycle.                                |
| 6     | Tag Release                               | After all documentation is updated and changes are committed, tag the release in Git.                                                        |

## Submitting a Pull Request (PR)

1.  **Push:** Push your feature or fix branch to your fork on GitHub:
    ```bash
    git push origin feature/your-descriptive-feature-name
    ```
2.  **Open PR:** Go to the original `dui-go` repository on GitHub (e.g., `https://github.com/duizendstra/dui-go`). GitHub should automatically detect your pushed branch and prompt you to create a Pull Request against the `main` branch (or the designated primary development branch).
3.  **Describe:**
    *   Fill out the Pull Request template comprehensively (if one exists).
    *   Clearly describe the problem your PR solves or the feature it adds.
    *   Explain the changes you've made.
    *   Link to any relevant GitHub issues (e.g., "Fixes #123" or "Implements feature discussed in #456").
4.  **Review:**
    *   Be responsive to feedback and code review comments.
    *   Ensure all CI checks (tests, linting) pass.
    *   The maintainers will review your PR. If it meets the project's standards and goals, it will be merged.

## Finding Ways to Contribute

*   Check the [GitHub Issues tab](https://github.com/duizendstra/dui-go/issues) (replace with your actual repo URL) for this repository. Look for issues tagged with `help wanted` or `good first issue`.
*   Review the [ROADMAP.MD](ROADMAP.MD) for planned features you might be interested in tackling.
*   Propose new features, enhancements, or identify bugs by creating a new issue.
*   Improving documentation or adding more test cases are always valuable contributions.

---
Thank you for contributing to `dui-go`!