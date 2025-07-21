# The Factory

## 1. Core Philosophy: The Developer Cockpit

Welcome to the `factory/` directory. This is the complete, self-contained automation framework for the Orion project. Its primary goal is to **reduce cognitive load** by providing a simple, safe, and powerful "cockpit" for the developer.

A developer should be able to focus on their intent (e.g., "I want to test my code") rather than the specific flags and commands required to do so. The `task` menu is the public API for this factory; everything else is an implementation detail contained within this directory.

## 2. Architecture: The "Menu / Workflow / Action" Pattern

The factory is built on a strict three-layer architectural pattern to ensure that logic is organized, maintainable, and easy to trace.

1.  **The Menu (`/Taskfile.yml`):**
    *   **What:** The single, user-facing entry point for all developers.
    *   **Why:** Provides a clean, high-level "public API" of commands (e.g., `task test`). It is a simple proxy that delegates every command to the factory. This is the only file a typical user needs to be aware of.

2.  **The Workflow (`factory/tasks/*.yml`):**
    *   **What:** The control panel that wires a menu command to its implementation.
    *   **Why:** This layer bridges the user's intent with the underlying logic. A file like `_test.yml` defines the task that calls the appropriate script. The `_` prefix in the filename prevents these workflow tasks from polluting the main `task --list` menu.

3.  **The Action (`factory/scripts/*.sh`):**
    *   **What:** The shell scripts that contain the core implementation logic.
    *   **Why:** This is where the real work happens. These scripts are enhanced with `gum` for a rich, interactive experience and contain all the necessary safety checks and commands.

## 3. Guiding Principle: Self-Contained Actions

A key architectural principle of this factory is that **individual action scripts should be self-contained**.

For example, the `release.sh` script contains all the logic it needs to run, including its own internal steps for testing and building. It does **not** call `task test` or `task build`.

**Rationale:**

This approach is a deliberate trade-off that prioritizes readability and portability over strict adherence to the DRY (Don't Repeat Yourself) principle at the script level.

*   **High Readability:** A developer can read a single script like `release.sh` and understand the entire end-to-end process without needing to trace calls across multiple other files.
*   **Portability & Future-Proofing:** Each script is a portable, self-contained unit of logic. This makes it significantly easier to "eject" this logic and migrate it into a future standalone CLI tool, which is a long-term goal of this framework.

## 4. How to Add a New Command

To add a new command (e.g., `task new-command`) to the framework, follow these five steps:

**Step 1: Create the Action Script**
Create a new, executable shell script in `factory/scripts/new_command.sh`. Place any reusable logic in the central `_automation_helpers.sh` library.

**Step 2: Create the Workflow Task**
Create a new YAML file in `factory/tasks/_new_command.yml`. This file will define the task that calls your new script, ensuring you pass arguments correctly.

\`\`\`yaml
# In factory/tasks/_new_command.yml
version: '3'
tasks:
  new-command:
    desc: "A description of your new command."
    silent: true
    cmds:
      - ./factory/scripts/new_command.sh {{.CLI_ARGS}}
\`\`\`

**Step 3: Include the Workflow**
Open the root `Taskfile.yml` and add an `include` for your new task file. Use a unique namespace to prevent collisions.

\`\`\`yaml
# In /Taskfile.yml
includes:
  _new_command: ./factory/tasks/_new_command.yml
  # ... other includes
\`\`\`

**Step 4: Expose the Command in the Menu**
In the root `Taskfile.yml`, create the user-facing facade task that delegates to your new workflow task.

\`\`\`yaml
# In /Taskfile.yml
tasks:
  new-command:
    desc: "A user-facing description of your command."
    cmds:
      - task: _new_command:new-command
        vars: { CLI_ARGS: '{{.CLI_ARGS}}' }
  # ... other tasks
\`\`\`

**Step 5: Document the Command in the Cockpit**
Finally, add your new command to the `default` task's help menu in the root `Taskfile.yml` so it is discoverable by users.