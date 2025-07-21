Excellent point, Kernel. A project of this complexity requires clear documentation at every level. The `tasks/` directory is the crucial link between the user-facing `Taskfile.yml` and the implementation logic in `scripts/`, so it absolutely deserves its own `README.md`.

This file will complete our documentation hierarchy, ensuring that any developer can understand the framework's design from top to bottom.

Here is a `README.md` crafted for the `tasks/` directory. It is designed to work in concert with the root and `scripts/` READMEs.

---

### `tasks/README.md`

```markdown
# Tasks Directory

## 1. Purpose

This directory contains the individual `Taskfile` configurations that are included by the root `Taskfile.yml`. Each file in this directory corresponds to a high-level command available to the user (e.g., `_commit.yml` implements the `task commit` command).

## 2. Role in the Automation Framework

Our framework follows a "Menu / Workflow / Action" design pattern to ensure clarity and separation of concerns. This directory represents the **"Workflow"** layer.

*   **Menu (`/Taskfile.yml`):** The user-facing command menu. It defines *what* commands are available. It acts as a clean, high-level facade.
*   **Workflow (`tasks/*.yml`):** This directory. It defines *how* a menu command is connected to its implementation. It's the bridge between the user's intent and the underlying logic.
*   **Action (`scripts/*.sh`):** The implementation layer. The shell scripts contain the actual logic that performs the work.

This structure keeps the root `Taskfile.yml` simple and readable, while the files in this directory provide the necessary wiring to the scripts.

## 3. File Naming and Structure

*   **Naming Convention:** All files must be prefixed with an underscore (e.g., `_commit.yml`, `_context.yml`). This is a standard `Task` convention that prevents the tasks within these files from appearing in the main `task --list` output, keeping the user-facing menu clean.
*   **Structure:** Each file is a standard `Taskfile` that defines one or more tasks. Typically, a file like `_commit.yml` will define a single task named `commit`.

## 4. Passing Arguments to Scripts

A critical responsibility of the tasks defined in this directory is to forward any command-line arguments to the underlying scripts. We use a standard pattern for this.

**The Pattern:**
```yaml
# tasks/_commit.yml
version: '3'

tasks:
  commit:
    desc: "(Action) Stages all changes and then commits them via a script."
    silent: true
    cmds:
      # Call the script, passing all command-line arguments through.
      - ./scripts/commit.sh {{.CLI_ARGS}}
```

The `{{.CLI_ARGS}}` variable is a feature of Go Task that captures all arguments passed after a `--` on the command line. This allows a user to run `task commit -- -m "My commit message"` and have the `-m "..."` part passed directly to `scripts/commit.sh`.

## 5. How to Add a New Task

To add a new command to the framework:

1.  **Create the Script:** Implement the core logic in a new shell script in the `scripts/` directory.
2.  **Create the Task File:** Create a new file in this directory, named `_your_task_name.yml`.
3.  **Define the Task:** Inside your new `.yml` file, define the task that calls your script, ensuring you pass `{{.CLI_ARGS}}`.
4.  **Include the Task File:** Open the root `Taskfile.yml` and add your new file to the `includes:` block.
5.  **Create the Facade:** In the root `Taskfile.yml`, add the user-facing facade task that delegates to your new task (e.g., `your-task: task: _your_task_name:your-task`).
6.  **Update the Help Menu:** Add your new command to the `default` task's help menu in the root `Taskfile.yml` so it is discoverable by users.
```
