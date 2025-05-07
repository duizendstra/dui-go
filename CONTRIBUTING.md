# Contributing to dui-go

Thank you for considering contributing to `dui-go`! We welcome improvements, bug fixes, and new features that help create a robust and useful Go utility library.

## Community Guidelines

Act professionally and respectfully. Be kind, considerate, and welcoming. Harassment or exclusionary behavior will not be tolerated.

## Getting Started

1.  **Prerequisites:** Ensure you have Go (`1.24` or later recommended) and Git installed and configured correctly.
2.  **Fork & Clone:** Fork the repository on GitHub (`github.com/duizendstra/dui-go` - *Ensure this is your correct public repository URL*) and clone your fork locally:
    ```bash
    # Replace YOUR_USERNAME with your actual GitHub username
    git clone https://github.com/YOUR_USERNAME/dui-go.git dui-go-fork
    cd dui-go-fork
    ```
3.  **Verify Setup:** Ensure you can build the library components and run tests:
    ```bash
    # Check if all packages build
    go build ./...
    # Run all tests
    go test ./...
    ```

## Making Changes

1.  **Create a Branch:** Before making changes, create a new branch from the `main` branch (or your primary development branch):
    ```bash
    git checkout main
    git pull origin main # Ensure your main is up-to-date
    git checkout -b feature/your-feature-name # Example: feature/add-new-cache-adapter
    # or
    git checkout -b fix/issue-description # Example: fix/improve-error-wrapping
    ```
2.  **Implement:** Make your code changes. Keep changes focused on a single feature or bug fix per branch.
3.  **Follow Style:**
    *   Adhere to standard Go formatting (`gofmt` or `goimports`).
    *   Run `go vet ./...` to catch common issues.
    *   If applicable, follow any linting guidelines established for the project (e.g., via a `Taskfile` if restored, or a specific linter).
4.  **Test:**
    *   **Automated:** If adding new functions or modifying existing ones, please add or update corresponding unit tests (`_test.go` files). Contributions to increase overall test coverage are highly encouraged. Run Go unit tests using:
        ```bash
        go test ./...
        # Consider running with the race detector as well
        go test -race ./...
        ```
    *   **Manual (if applicable):** If your changes affect observable behavior in a way not easily unit-tested, describe the manual steps you took to verify.
5.  **Commit:** Commit your changes using clear and descriptive commit messages. Consider following the [Conventional Commits](https://www.conventionalcommits.org/) specification (e.g., `feat: ...`, `fix: ...`, `refactor: ...`, `docs: ...`).
    ```bash
    git add .
    git commit -m "feat(cache): Add Redis adapter for Cache interface"
    ```

## Submitting a Pull Request

1.  **Push:** Push your feature or fix branch to your fork on GitHub:
    ```bash
    git push origin feature/your-feature-name
    ```
2.  **Open PR:** Go to the original `github.com/duizendstra/dui-go` repository on GitHub. GitHub should automatically detect your pushed branch and prompt you to create a Pull Request against the `main` branch (or your primary development branch).
3.  **Describe:** Fill out the Pull Request template (if one exists), clearly describing the problem you're solving and the changes you've made. Link to any relevant GitHub issues (e.g., "Fixes #123" or "Relates to #456").
4.  **Review:** Respond to any feedback or code review comments. The maintainers will review your PR and merge it if it meets the project's standards.

## Finding Ways to Contribute

Check the [GitHub Issues tab](https://github.com/duizendstra/dui-go/issues) for this repository. Look for issues tagged with `help wanted` or `good first issue`, or feel free to pick up any unassigned issue that interests you. You can also propose new features or enhancements by creating a new issue.

Thank you for contributing to `dui-go`!