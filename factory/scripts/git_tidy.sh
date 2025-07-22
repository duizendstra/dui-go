#!/bin/bash
# factory/scripts/git_tidy.sh
#
# WHAT: A workflow-driven script to help developers with Git branch hygiene.
#       It provides interactive tools for finishing merged branches and pruning
#       stale local branches.

set -e

# --- Configuration & Setup ---
# shellcheck source=./_automation_helpers.sh
source "$(dirname "$0")/_automation_helpers.sh"

ORIGINAL_BRANCH=$(git rev-parse --abbrev-ref HEAD)
STASH_PERFORMED="false"

# --- Safety Net ---
cleanup() {
  if [[ "$(git rev-parse --abbrev-ref HEAD)" != "$ORIGINAL_BRANCH" ]]; then
    echo
    # FIX: Added '--' to fix gum style flag error.
    gum style --bold -- "--> Returning to original branch '$ORIGINAL_BRANCH'..."
    git checkout "$ORIGINAL_BRANCH" >/dev/null 2>&1
  fi
  if [ "$STASH_PERFORMED" = "true" ]; then
    echo
    # FIX: Added '--' to fix gum style flag error.
    gum style --bold -- "--> Re-applying stashed changes..."
    git stash pop || echo "⚠️ Warning: Could not auto-apply stashed changes due to a conflict."
  fi
}
trap cleanup EXIT

# --- ACTION FUNCTIONS ---

finish_merged_branch() {
    gum style --bold --padding "1 0" -- "--- Finish Merged Branch Workflow ---"
    
    if [[ "$ORIGINAL_BRANCH" == "$MAIN_BRANCH" ]]; then
        gum style --foreground 9 "Error: You are already on the main branch. There is no branch to finish."
        exit 1
    fi
    
    if ! gum confirm "This will delete your local branch '$ORIGINAL_BRANCH' and switch to '$MAIN_BRANCH'. Are you sure it has been merged?"; then
        echo "Aborted by user."
        exit 0
    fi
    
    echo "--> Switching to '$MAIN_BRANCH'..."
    git checkout "$MAIN_BRANCH"
    
    echo "--> Pulling latest changes..."
    git pull origin "$MAIN_BRANCH"
    
    echo "--> Deleting local branch '$ORIGINAL_BRANCH'..."
    git branch -d "$ORIGINAL_BRANCH"
    
    gum style --bold "✅ Successfully cleaned up '$ORIGINAL_BRANCH' and updated '$MAIN_BRANCH'."
    if gum confirm "Also prune other stale branches now?"; then
        prune_stale_branches
    fi
}

prune_stale_branches() {
    gum style --bold --padding "1 0" -- "--- Prune Stale Branches ---"
    echo "--> Fetching remote state and searching for stale branches..."
    git fetch --all --prune
    
    local temp_main_branch_checkout=false
    if [[ "$(git rev-parse --abbrev-ref HEAD)" != "$MAIN_BRANCH" ]]; then
        temp_main_branch_checkout=true
        git checkout "$MAIN_BRANCH" >/dev/null 2>&1
    fi
    
    local remote_branches
    remote_branches=$(git branch -r | sed 's|origin/||' | sed 's/^[ \t]*//')
    local merged_local_branches
    merged_local_branches=$(git branch --merged main | grep -vE '^\*|main$' | sed 's/^[ \t]*//')
    local stale_branches
    stale_branches=$(comm -23 <(echo "$merged_local_branches" | sort) <(echo "$remote_branches" | sort))
    
    if [ "$temp_main_branch_checkout" = true ]; then
        git checkout "$ORIGINAL_BRANCH" >/dev/null 2>&1
    fi

    if [ -z "$stale_branches" ]; then
        gum style --bold "✅ No stale local branches found."
        return
    fi
    
    gum style --bold "The following stale branches can be safely deleted:"
    echo "$stale_branches" | gum style --faint

    if gum confirm "Proceed with deletion?"; then
        echo "$stale_branches" | xargs git branch -d
        gum style --bold "✅ Stale branches deleted."
    else
        gum style --bold "Aborted. No branches were deleted."
    fi
}

# --- Main Script Logic ---

if ! prompt_to_stash_if_dirty; then
    exit 1
fi

gum style --bold --padding "1 0" -- "🌿 What Git hygiene task would you like to perform?"
CHOICE=$(gum choose     "Finish Merged Branch (delete current branch, go to main)"     "Prune Old Branches (stay on current branch)"     "Quit")

case "$CHOICE" in
    "Finish Merged Branch (delete current branch, go to main)")
        finish_merged_branch
        ;;
    "Prune Old Branches (stay on current branch)")
        prune_stale_branches
        ;;
    *)
        gum style --bold "Aborted."
        exit 0
        ;;
esac

echo
gum style --bold "✅ Operation complete."
