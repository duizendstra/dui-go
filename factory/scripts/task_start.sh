#!/bin/bash
# factory/scripts/start_task.sh
#
# WHY:  Handles starting a new feature branch safely.
# WHAT: Supports both interactive and parameterized branch creation.
#       - Interactive: `task task-start`
#       - Parameterized: `task task-start -- <type> <scope> <description>`

set -e

# --- Phase 1: State Verification ---
STASH_PERFORMED=false
# Check for uncommitted changes, but only if not running in a CI environment
if [ -z "$CI" ] && ! git diff --quiet --exit-code; then
  gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "⚠️ You have uncommitted changes."
  if gum confirm "Stash them and bring them to the new branch?"; then
    git stash
    STASH_PERFORMED=true
    echo "✅ Changes stashed."
  else
    echo "Aborted by user. Please commit or stash your changes."
    exit 1
  fi
fi

# --- Phase 2: Branch Creation (Interactive or Parameterized) ---
BRANCH_TYPE=$1
PBI_ID=$2
DESCRIPTION=$3

# If arguments are not provided, enter interactive mode.
if [ -z "$BRANCH_TYPE" ] || [ -z "$PBI_ID" ] || [ -z "$DESCRIPTION" ]; then
  gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "🌿 Let's create a new branch."

  echo "Select a branch type:"
  BRANCH_TYPE=$(gum choose "feature" "fix" "docs" "style" "refactor" "test" "chore" "factory")
  if [ -z "$BRANCH_TYPE" ]; then exit 1; fi

  echo "Enter the PBI number or scope (e.g., SFB-003):"
  PBI_ID=$(gum input --placeholder "SFB-XXX")
  if [ -z "$PBI_ID" ]; then exit 1; fi

  echo "Enter a short description (use-kebab-case):"
  DESCRIPTION=$(gum input --placeholder "implement-new-feature")
  if [ -z "$DESCRIPTION" ]; then exit 1; fi
fi

BRANCH_NAME="${BRANCH_TYPE}/${PBI_ID}/${DESCRIPTION}"
gum confirm "Create branch '$BRANCH_NAME'?" || exit 0

# --- Phase 3: Git Operation ---
gum spin --spinner dot --title "Creating branch..." -- git checkout -b "$BRANCH_NAME"

# --- Phase 4: Restore Stashed Changes (if any) ---
if [ "$STASH_PERFORMED" = true ]; then
  echo "--> Re-applying your stashed changes..."
  git stash pop
  echo "✅ Your work has been restored."
fi

gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "✅ Success! You are now on branch '$BRANCH_NAME'."