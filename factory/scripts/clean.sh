#!/bin/bash
# factory/scripts/clean.sh
#
# WHAT: A script to clean up local project files and caches.
# WHY:  Provides a safe, interactive utility for removing temporary files
#       and build artifacts to free up space and ensure a clean state.

set -e

# --- Configuration & Setup ---
# shellcheck source=./_automation_helpers.sh
source "$(dirname "$0")/_automation_helpers.sh"

# --- ACTION FUNCTIONS ---

clean_project_files_only() {
    gum style --bold --padding "1 0" -- "--- Clean Local Project Files ---"
    echo "--> Removing compiled binaries (./bin)..."
    rm -rf ./bin
    echo "--> Cleaning Go build cache..."
    go clean
    echo "--> Removing temporary context files..."
    rm -f context*
    echo "--> Removing Task runner cache (./.task)..."
    rm -rf ./.task
    gum style --bold "✅ Project files cleaned."
}

clean_full_system() {
    gum style --bold --padding "1 0" -- "--- Full System Clean ---"
    
    if ! gum confirm "DANGER: This will purge Go module caches and ALL unused Docker images, containers, and volumes on your system. This can affect other projects. Proceed?"; then
        echo "Aborted by user."
        exit 0
    fi

    clean_project_files_only
    
    echo "--> Purging Go module and test caches..."
    go clean -modcache -testcache

    echo "--> Pruning all unused Docker resources..."
    docker system prune -af --volumes

    gum style --bold "✅ Full system clean complete."
}


# --- Main Script Logic ---

# Although this script doesn't modify git history, it's good practice to
# check for uncommitted changes before deleting files, as it could be confusing.
if ! prompt_to_stash_if_dirty; then
    exit 1
fi

gum style --bold --padding "1 0" -- "🧹 What would you like to clean?"
CHOICE=$(gum choose \
    "Clean Local Project Files Only" \
    "Full System Clean (Project Files + Docker & Go Caches)" \
    "Quit")

case "$CHOICE" in
    "Clean Local Project Files Only")
        clean_project_files_only
        ;;
    "Full System Clean (Project Files + Docker & Go Caches)")
        clean_full_system
        ;;
    *)
        gum style --bold "Aborted."
        exit 0
        ;;
esac

echo
gum style --bold "✅ Operation complete."
