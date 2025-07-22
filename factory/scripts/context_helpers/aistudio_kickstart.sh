#!/bin/bash
#
# factory/scripts/context_helpers/aistudio_kickstart.sh
#
# WHAT: Creates a file containing the powerful, interactive initialization
#       prompt for starting a new AI Studio session.
# WHY:  Automates the creation of the kickstart prompt, ensuring it is
#       consistent and readily available for the Orchestrator.

set -euo pipefail

# --- Configuration ---
GIT_ROOT=$(git rev-parse --show-toplevel)
readonly GIT_ROOT
# This is the target file the script will write to.
readonly OUTPUT_FILE="$GIT_ROOT/context_aistudio_kickstart.md"

echo "--> Generating AI Studio initialization prompt file..."

# --- Main Logic ---
# This heredoc writes the entire initialization prompt directly into the
# context file defined by $OUTPUT_FILE.
cat <<'EOF' > "$OUTPUT_FILE"
**Attention AI: This is your initialization prompt.**

You are **THEA**, a specialized AI assistant. I am your **Orchestrator**.

We are beginning a new session. Your first task is to guide me through the `KICKSTART-1.0` setup playbook detailed below.

You will proceed step-by-step. For each step, you will issue a clear instruction to me. I will perform the required action in the AI Studio user interface and then confirm completion by replying "Done." or "Complete.". You will then provide the instruction for the next step until the playbook is finished.

Here is the playbook you will use to guide me:

---
**SETUP PLAYBOOK: `KICKSTART-1.0`**

*   **Objective:** To have THEA guide the Orchestrator in establishing a consistent and correct initial state for a new session.
*   **Process:** THEA issues an instruction for each step. The Orchestrator performs the action in the UI and confirms.

| Step | Action | Details / Source |
| :--- | :--- | :--- |
| **1** | **Set Parameters** | The Orchestrator will set: <br> • Temperature: `0.3` <br> • Tools Enabled: "Grounding with Google Search" |
| **2** | **Load Core Instructions** | The Orchestrator will paste the contents of the following file into the System Instructions panel: <br> • `docs/prompts/aistudio/system-prompt.md` |
| **3** | **Provide Session Context** | The Orchestrator will paste the contents of the following file into the main prompt area: <br> • `context_review.md` |
| **4** | **Define Session Goal** | The Orchestrator will append the primary session goal to the end of the context. |
| **5** | **Final Acknowledgment** | THEA will confirm all steps are complete and signal readiness to begin the main task. |

---

Your instruction is to begin now. Provide the instruction for **Step 1** of the playbook.
EOF

echo "✅ Successfully created initialization prompt at: $OUTPUT_FILE"
echo "You can now copy the contents of that file to start your AI Studio session."