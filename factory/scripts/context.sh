#!/bin/bash
#
# Generates a comprehensive context file for the project by analyzing the Git
# state and the Go modules defined in the go.work file.
#

set -e # Exit immediately if a command exits with a non-zero status.

# --- Helper Functions ---

# Prints a styled header for a section of the context file.
# Usage: print_header "My Section Title"
print_header() {
  echo
  # Uses gum for consistent, styled output.
  gum style --border double --padding "0 1" --border-foreground 212 "$1"
  echo
}

# Discovers all Go modules from the go.work file and prints a detailed
# summary of each one, including its files and direct dependencies.
print_go_modules_context() {
  print_header "Go Workspace Modules"
  
  # Use 'go work edit -json' as the single source of truth for module paths.
  # This is the most robust way to discover modules in a workspace.
  local modules
  modules=$(go work edit -json | jq -r '.Use[].DiskPath')

  if [ -z "$modules" ]; then
    echo "No modules found in go.work file."
    return
  fi

  # Iterate through each discovered module directory.
  for module_dir in $modules; do
    # Run commands for each module in a subshell to prevent 'cd' from
    # affecting the main script's working directory.
    (
      cd "$module_dir" || continue
      
      # Get the full module path (e.g., github.com/duizendstra/dui-go/gcs) from its go.mod file.
      local module_name
      module_name=$(go mod edit -json | jq -r '.Module.Path')
      
      # *** FIXED LINE ***
      # Added '--' to prevent gum from interpreting the styled text as a flag.
      gum style --bold -- "--- Module: $module_name ---"
      echo "Location: ./$module_dir"
      echo
      
      echo "Files:"
      # Find all .go files, excluding tests for brevity. The sed command cleans up the './' prefix.
      find . -name "*.go" ! -name "*_test.go" -print | sed 's|^\./||' | gum format
      echo

      echo "Direct Dependencies:"
      # Parse the go.mod file to list only direct dependencies, which is more
      # informative for context than listing all transitive dependencies.
      go mod edit -json | jq -r '.Require[] | select(.Indirect | not) | "  - \(.Path)@\(.Version)"' | gum format
      echo
    )
  done
}

# Generates the default context, which is tailored for writing a commit message.
# It focuses on what has changed.
generate_commit_context() {
  local output_file="context_commit.md"
  # Group all command outputs and redirect them into the final file.
  {
    print_header "Git Status"
    git status --short
    echo

    print_header "Staged Git Diff"
    # Provide a statistical summary and the full patch for detailed review.
    git diff --staged --stat --color=always
    echo
    git diff --staged --patch --color=always
    echo

    # Include the module summary to see how the code structure is affected.
    print_go_modules_context

  } > "$output_file"
  
  echo "✅ Commit context generated at: $output_file"
}

# Generates a comprehensive export of the entire project's structure and
# configuration, intended for a full project analysis.
generate_export_all_context() {
  local output_file="contextvibes_export_project.md"
  {
    print_header "Project README"
    cat README.md
    echo
    
    print_header "Project Structure"
    # Use 'ls -R' for a recursive, detailed file listing.
    ls -R
    echo

    # Include the full module summary.
    print_go_modules_context

  } > "$output_file"
  echo "✅ Full project context exported to: $output_file"
}


# --- Main Execution ---
main() {
  # Default to the 'commit' context if no specific type is requested.
  local context_type="${1:---default}"

  case "$context_type" in
    --export-all)
      generate_export_all_context
      ;;
    *) # Handles the default case or any other arguments.
      generate_commit_context
      ;;
  esac
}

# Pass all script arguments to the main function.
main "$@"
