# Scripts Directory

## 1. Purpose

This directory contains all the shell scripts that implement the core logic for the project's automation framework.

These scripts represent the "Action" layer in our "Menu / Workflow / Action" automation philosophy. They are designed to be called by the `Taskfile.yml` configurations located in the `tasks/` directory and should contain the actual implementation details for a given command.

## 2. Design Philosophy

*   **Separation of Concerns:** The root `Taskfile.yml` is the user-facing "Menu." The `tasks/*.yml` files define the "Workflow." These scripts are the "Action"—the "how." This separation keeps the `Taskfile` clean and makes the underlying logic easier to find and maintain.
*   **Single Responsibility:** Each script should, as much as possible, be responsible for a single, well-defined task (e.g., `commit.sh` handles committing, `start_task.sh` handles starting a new task).
*   **Don't Repeat Yourself (DRY):** All generic, reusable logic that could be used by more than one script **must** be placed in the central helper library.

## 3. The Automation Helper Library

The most important file in this directory is `_automation_helpers.sh`.

*   **What It Is:** A library of shared variables and bash functions that are used across multiple scripts. The leading underscore (`_`) ensures it is listed first and signals its role as an internal library.
*   **Purpose:** It is the **single source of truth** for common logic. This prevents code duplication and ensures that core behaviors (like checking for uncommitted changes or validating branch names) are consistent everywhere.
*   **Usage:** To use the helper functions in a new script, add the following line at the top:
    ```sh
    source "$(dirname "$0")/_automation_helpers.sh"
    ```

**Golden Rule:** Before writing any new code in a script, ask yourself: "Is this logic generic enough to be used by another script?" If the answer is yes, it belongs in `_automation_helpers.sh`.

## 4. Naming Conventions

Scripts should be named clearly and consistently to reflect their purpose. We follow a `verb_noun.sh` or `context_verb.sh` pattern.

*   **Examples:**
    *   `commit.sh` (Action script)
    *   `start_task.sh` (Action script)
    *   `context_export_automation.sh` (Context generation script)

## 5. Contributing a New Script

When adding a new script to the framework, follow these steps:

1.  **Check the Helper:** Review `_automation_helpers.sh` to see if any functions you need already exist.
2.  **Factor Out Logic:** As you write your script, identify any new, generic logic and add it to `_automation_helpers.sh` first.
3.  **Source the Helper:** Add `source "$(dirname "$0")/_automation_helpers.sh"` to the top of your new script.
4.  **Write the Script:** Implement the specific logic for your new task, calling helper functions where appropriate.
5.  **Make it Executable:** Run `chmod +x your_new_script.sh` to ensure it can be executed.
6.  **Update Taskfile:** Add a corresponding task in the `tasks/` directory that calls your new script.
