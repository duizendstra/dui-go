# To learn more about how to use Nix to configure your environment
# see: https://firebase.google.com/docs/studio/customize-workspace
{ pkgs, ... }: {
  # Which nixpkgs channel to use (stable recommended for reproducibility).
  channel = "stable-24.11";

  # Packages available in the development environment.
  # Use https://search.nixos.org/packages to find more.
  packages = [
    # --- Language Toolchain ---
    pkgs.go

    # --- Version Control ---
    pkgs.git

    # --- Utilities ---
    # pkgs.patch # Optional: Standard patching utility, remove if not used.
    pkgs.jq    # Command-line JSON processor.
    pkgs.tree  # Directory structure viewer.
  ];

  # Environment variables specific to this workspace.
  env = { };

  # Project IDX specific settings.
  idx = {
    # VS Code extensions to install. Find IDs at https://open-vsx.org/
    extensions = [
      "golang.go"
      "GitHub.vscode-pull-request-github"
    ];

    # Workspace previews configuration (disabled for this library).
    previews = {
      enable = false;
      previews = { };
    };

    # Workspace lifecycle hooks (empty for this library).
    workspace = {
      # Runs once when the workspace is created.
      onCreate = {
        # Example: download-deps = "go mod download";
      };
      # Runs every time the workspace starts.
      onStart = { };
    };
  };
}