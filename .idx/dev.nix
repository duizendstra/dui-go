# This file defines the development environment for a Firebase Studio workspace using Nix.
# It specifies all the packages, services, and editor extensions needed for the project.
# This consolidated version removes duplicates and organizes components for better clarity.
{ pkgs, ... }:

let
  # Imports the custom Nix derivation for the 'contextvibes' tool from a separate file.
  # This declarative approach makes the tool a first-class package in the environment,
  # which is more reproducible and cacheable than installing it with an imperative script.
  contextvibes = import ./contextvibes.nix { pkgs = pkgs; };
in
{
  # Specifies the Nixpkgs channel to use. Pinning to a specific channel like "stable-25.05"
  # ensures that everyone on the team gets the exact same package versions, leading to a
  # highly reproducible environment.
  channel = "stable-25.05";

  # Defines the packages to be installed in the development environment.
  # Packages are de-duplicated and grouped by function for better readability.
  packages = with pkgs; [
    # --- Core Go Development Toolchain ---
    go # The Go compiler and core toolchain.
    gopls # The official Go Language Server for IDE features.
    gotools # A suite of supplementary Go tools used by IDE extensions.
    delve # The Go debugger, essential for step-debugging.
    goimports-reviser # A tool to automatically format and revise Go import statements.
    golangci-lint # A fast Go linter that runs multiple linters in parallel.
    gcc # The GNU Compiler Collection, required by Go for cgo support.

    # --- Automation, Containers & Cloud ---
    go-task # A task runner for project automation (see Taskfile.yml).
    docker-compose # For orchestrating local multi-container Docker applications.
    google-cloud-sdk # The `gcloud` CLI for interacting with Google Cloud Platform.

    # --- Code Quality & Formatting (Non-Go) ---
    shellcheck # Linter for finding bugs in shell scripts.
    shfmt # Auto-formatter for shell scripts.
    nodejs # JavaScript runtime, required for markdownlint-cli.
    nodePackages.markdownlint-cli # Linter to enforce standards in Markdown files.

    # --- Version Control & CLI Utilities ---
    git # The version control system for managing source code.
    gh # The official GitHub CLI for interacting with GitHub.
    jq # A command-line JSON processor for scripting.
    yq-go # A portable command-line YAML processor.
    tree # A utility to display directory structures as a tree.
    file # A utility to determine file types.
    gum # A tool for creating glamorous, interactive shell scripts.

    # --- Custom Project Tools ---
    contextvibes # The custom-built 'contextvibes' CLI tool, managed by its own Nix file.
  ];

  # Enables and starts system-level services within the environment.
  services.docker.enable = true; # Starts the Docker daemon, required for building and running containers.

  # Configures the Firebase Studio workspace environment.
  idx = {
    # Specifies VS Code extensions to install automatically.
    # This ensures a consistent and fully-featured editor experience for all developers.
    extensions = [
      # --- Core Language Support ---
      "golang.go" # Official Go extension (debugging, testing, linting).
      "ms-python.python" # Python language support.
      "ms-python.debugpy" # Python debugging support.

      # --- Tooling & Integration ---
      "ms-azuretools.vscode-docker" # Docker integration and container management.
      "task.vscode-task" # Adds support for Go Task ('Taskfile.yml').

      # --- Code Quality ---
      "DavidAnson.vscode-markdownlint" # Integrates markdownlint into the editor.
      "timonwong.shellcheck" # Integrates shellcheck for live linting of shell scripts.

      # --- Version Control ---
      "GitHub.vscode-pull-request-github" # GitHub Pull Request and Issues integration.
      "eamodio.gitlens" # Supercharges the Git capabilities built into VS Code.
    ];
  };

  # Note on Omitted Configurations:
  #
  # - `idx.previews`: This section is omitted because its default value (`enable = false`)
  #   is the desired behavior for this command-line application.
  #
  # - `idx.workspace.onCreate` & `idx.workspace.onStart`: These lifecycle hooks are not
  #   defined here. Custom tools are included declaratively in the `packages` list, which
  #   is the preferred Nix approach for reproducibility and caching over setup scripts.
  #   These hooks can be added if other non-package setup scripts are needed.
}
