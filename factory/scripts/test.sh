#!/bin/bash
#
# Runs the test suite for every Go module in the workspace.
# This script is module-aware and aggregates coverage reports.
#

set -e # Exit immediately if a command exits with a non-zero status.

# --- Main Execution ---
main() {
  echo "Discovering modules in go.work..."
  # Use the go.work file as the single source of truth for module locations.
  local modules
  modules=$(go work edit -json | jq -r '.Use[].DiskPath')

  if [ -z "$modules" ]; then
    echo "❌ No modules found in go.work file. Nothing to test."
    exit 1
  fi

  # Prepare a root coverage file in a temporary location.
  local root_coverage_file
  root_coverage_file=$(mktemp)
  echo "mode: set" > "$root_coverage_file"

  # Loop through each module to run tests and collect coverage.
  for module_dir in $modules; do
    echo
    echo "--- Testing Module: $module_dir ---"
    # Run tests in a subshell to isolate the 'cd' command.
    (
      cd "$module_dir"
      # Run tests with race detector, verbosity, and generate a coverage profile.
      go test -race -v -coverprofile=coverage.out -covermode=atomic ./...
    )

    # Append the module's coverage data to the root coverage file, skipping the first line ("mode: set").
    local module_coverage_file="$module_dir/coverage.out"
    if [ -f "$module_coverage_file" ]; then
      tail -n +2 "$module_coverage_file" >> "$root_coverage_file"
      rm "$module_coverage_file" # Clean up the individual report
    fi
  done

  echo
  echo "--- ✅ All Module Tests Passed ---"
  echo
  echo "--- Total Project Coverage ---"
  # Calculate and display the total functional coverage from the aggregated report.
  go tool cover -func="$root_coverage_file"

  # Clean up the temporary aggregated file
  rm "$root_coverage_file"
}

# --- Run the main function ---
main
