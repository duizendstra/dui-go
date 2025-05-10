{ pkgs, ... }: {
  # Which nixpkgs channel to use. (https://status.nixos.org/)
  channel = "stable-24.11"; # Or choose a specific Nixpkgs commit/tag

  # Use https://search.nixos.org/packages to find packages for Go development
  packages = [
    # --- Core Go Development ---
    pkgs.go # The Go compiler and runtime
    pkgs.gotools # Includes gopls, goimports, delve, etc.

    # --- Version Control ---
    pkgs.git # Essential version control system
    pkgs.gh # Github CLI
  ];

  # Sets environment variables in the workspace
  env = { };

  idx = {
    # Search for extensions on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      # --- Go Language Support ---
      "golang.go" # Official Go extension (debugging, testing, linting/formatting)

      # --- Version Control ---
      "GitHub.vscode-pull-request-github" # GitHub Pull Request and Issues integration
    ];

    workspace = {
      # Runs when a workspace is first created with this `dev.nix` file
      onCreate = {
        # Script to install contextvibes CLI into ./bin
        installContextVibesCli = ''
          echo "Attempting to install contextvibes CLI into ./bin ..."

          if ! command -v go &> /dev/null # Checks if 'go' command is available
          then
              echo "Go command could not be found by 'command -v go', skipping contextvibes installation."
              # This path might be taken if Nix hasn't fully set up the PATH yet when onCreate runs
              # or if pkgs.go isn't providing the 'go' command as expected in this hook's environment.
              # However, pkgs.go *should* make it available.
          else
            # Ensure we are in the workspace root for consistent relative paths
            # cd "$WS_ROOT" || { echo "ERROR: Failed to cd to workspace root $WS_ROOT"; exit 1; }
            # "$(pwd)" within onCreate hook should be the workspace root.

            LOCAL_BIN_DIR="$(pwd)/bin" # Use pwd, should be workspace root
            mkdir -p "$LOCAL_BIN_DIR"
            echo "Target directory for contextvibes: $LOCAL_BIN_DIR"

            # Using export GOBIN might not persist if 'go install' forks a subshell
            # that doesn't inherit it, though often it does.
            # A safer bet is to specify the output directly if 'go install' supports -o with remote paths,
            # but it doesn't for remote packages. GOBIN is the standard way.
            export GOBIN="$LOCAL_BIN_DIR"
            echo "Running: GOBIN=$GOBIN go install github.com/contextvibes/cli/cmd/contextvibes@v0.1.1"

            if go install github.com/contextvibes/cli/cmd/contextvibes@v0.1.1; then
              echo "Successfully installed contextvibes to $LOCAL_BIN_DIR/contextvibes"
              echo "You can run it using: $LOCAL_BIN_DIR/contextvibes or ./bin/contextvibes"
              echo "For convenience, you might add $LOCAL_BIN_DIR to your PATH environment variable."
              echo "  Example for .bashrc/.zshrc: export PATH=\"$(pwd)/bin:\$PATH\""
              # chmod +x is often not needed as 'go install' usually sets execute permissions.
              # but it doesn't hurt to ensure it.
              chmod +x "$LOCAL_BIN_DIR/contextvibes" || echo "Note: chmod +x on contextvibes failed (permissions?)."
            else
              echo "ERROR: Failed to install contextvibes. Check Go environment and network."
              # Consider 'exit 1' here if contextvibes is critical for the dev workflow from the start.
            fi
            unset GOBIN # Clean up GOBIN
          fi
        '';
      };
      # Runs every time a workspace is started
      onStart = {
        checkContextVibes = ''
          if [ -f "$(pwd)/bin/contextvibes" ]; then
            echo "ContextVibes CLI found in ./bin. Version: $($(pwd)/bin/contextvibes version)"
          else
            echo "ContextVibes CLI not found in ./bin. It may need to be installed (usually via workspace.onCreate)."
          fi
        '';
      };
    };

    # Enable previews and customize configuration if you're running web services
    previews = {
      enable = false; # Correct for a library project like dui-go
    };
  };
}
