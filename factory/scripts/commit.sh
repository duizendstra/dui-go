#!/bin/bash
# factory/scripts/commit.sh
#
# WHAT: Provides a safe, interactive, and guided workflow for creating Git commits.
#       Also supports a non-interactive mode for AI-driven automation.
# WHY:  Prevents common errors and improves commit quality.

set -e

# --- Configuration & Setup ---
# shellcheck source=./_automation_helpers.sh
source "$(dirname "$0")/_automation_helpers.sh"

# --- Safety Checks ---
prevent_action_on_main_branch "commit"
validate_current_branch_compliance

# --- Main Logic: Non-Interactive vs. Interactive ---

# If arguments are passed (e.g., -m "..."), use them directly for a non-interactive commit.
if [ -n "$1" ]; then
    echo "--> Staging all changes for non-interactive commit..."
    git add .
    echo "--> Committing non-interactively with provided message..."
    git commit "$@"
    gum style --foreground 212 "✅ Commit successful."
    exit 0
fi

# --- Fallback to Interactive Mode if no arguments are provided ---

# --- Phase 1: Interactive Staging ---
echo "--> Analyzing changed files..."
GIT_STATUS_LINES=$(git status --porcelain)

if [ -z "$GIT_STATUS_LINES" ]; then
    gum style --bold "✅ No changes to commit. Working tree is clean."
    exit 0
fi

gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "SELECT FILES TO STAGE FOR COMMIT"
echo "Status: M=Modified, A=Added, D=Deleted, R=Renamed, ??=Untracked"
echo "Use [space] to select, [enter] to confirm."

SELECTED_LINES=$(echo "$GIT_STATUS_LINES" | gum filter --no-limit --placeholder="Choose files to stage for the commit...")

if [ -z "$SELECTED_LINES" ]; then
    echo "No files selected. Aborting commit."
    exit 0
fi

# --- Staging Logic ---
echo "--> Staging selected files..."

echo "$SELECTED_LINES" | grep -v '^>' | while IFS= read -r line; do
    if [ -z "$line" ]; then continue; fi
    path_to_add=$(echo "$line" | sed -e 's/^...//' -e 's/.* -> //')
    echo "Staging: $path_to_add"
    git add -- "$path_to_add"
done

gum style --faint "Staging complete."

# --- Phase 2: Interactive Commit Message ---
gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "COMPOSE COMMIT MESSAGE"

COMMIT_TYPE=$(gum choose "feat" "fix" "docs" "style" "refactor" "test" "chore" "factory")
COMMIT_SCOPE=$(gum input --placeholder="(optional scope, e.g., 'api', 'db')")
COMMIT_SUBJECT=$(gum input --placeholder="Short description (imperative mood, e.g., 'Add new endpoint')")
COMMIT_BODY=$(gum write --placeholder="Longer description (optional)...")

if [ -n "$COMMIT_SCOPE" ]; then
    COMMIT_SCOPE="($COMMIT_SCOPE)"
fi

COMMIT_HEADER="${COMMIT_TYPE}${COMMIT_SCOPE}: ${COMMIT_SUBJECT}"

# --- Phase 3: Final Confirmation and Commit ---
gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "FINAL CONFIRMATION"
gum style --bold "Commit with the following message?"
echo
gum style --background 236 --padding "1" "${COMMIT_HEADER}"$'\n\n'"${COMMIT_BODY}"
echo

if ! gum confirm "Proceed with commit?"; then
    echo "Aborted by user."
    echo "To unstage the files, run: git reset"
    exit 1
fi

echo "--> Committing..."
git commit -m "$COMMIT_HEADER" -m "$COMMIT_BODY"

gum style --foreground 212 "✅ Commit successful."
