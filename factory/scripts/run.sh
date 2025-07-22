#!/bin/bash
# factory/scripts/run.sh
# An interactive script to run the application in various modes.

set -e

# --- Helper Functions ---

# Function to build the local binary, used by one of the run modes.
build_binary() {
    gum style --bold "--> Building Go binary into ./bin/..."
    mkdir -p ./bin
    gum spin --spinner dot --title "Compiling..." -- \
        go build -v -o ./bin/ms-graph-worker ./src/cmd/main.go
}

# Function to build the local container image, used by container run modes.
build_container() {
    local image_fqn="${GCP_REGION}-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/${ARTIFACT_REGISTRY_NAME}/${CONTAINER_NAME}"
    local git_commit
    git_commit=$(git rev-parse --short HEAD)

    gum spin --spinner dot --title "Building local container image..." -- \
    docker buildx build \
        --platform linux/amd64 \
        -t "${image_fqn}:${git_commit}" \
        -t "${image_fqn}:latest" \
        --file ./factory/Dockerfile \
        .
}

# --- Main Logic ---

# Get the current gcloud user for display purposes.
CURRENT_USER=$(gcloud config get-value account --quiet 2>/dev/null || echo "gcloud-user-not-found")

# --- Phase 1: Select the Target Entity ---
gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 "--- 1. Select Data Entity to Export ---"
ENTITY_CHOICE=$(gum choose \
    "users" \
    "groups" \
    "group_members" \
    "user_licenses" \
    "subscribed_skus" \
    "audit_logs" \
    "Quit")

if [[ "$ENTITY_CHOICE" == "Quit" ]]; then
    echo "Aborted."
    exit 0
fi

export TARGET_ENTITY="$ENTITY_CHOICE"
echo "--> Target entity set to: $(gum style --bold "$TARGET_ENTITY")"
echo

# --- Phase 2: Select the Run Mode ---
gum style --border double --border-foreground 212 --padding "1 2" -- "--- 2. Select Application Runner Mode ---"
RUN_CHOICE=$(gum choose \
  "Run from local source code as (${CURRENT_USER})" \
  "Run from compiled binary as (${CURRENT_USER})" \
  "Run in a single container" \
  "Trigger a remote job on Cloud Run" \
  "Quit")

# Execute the chosen run mode.
case "$RUN_CHOICE" in
    "Run from local source code as (${CURRENT_USER})")
        echo "--> Running Go application from source (verbose debug mode)..."
        # For local source runs, we need to ensure USE_WORKLOAD_IDENTITY is false.
        export USE_WORKLOAD_IDENTITY="false"
        go run ./src/cmd/main.go
        ;;

    "Run from compiled binary as (${CURRENT_USER})")
        build_binary
        export USE_WORKLOAD_IDENTITY="false"
        gum spin --spinner dot --title "Running the locally built binary..." --show-output -- \
            ./bin/ms-graph-worker
        ;;

    "Run in a single container")
        build_container
        gum style --bold "--> Starting single container for entity '${TARGET_ENTITY}'..."
        # Pass GIT_COMMIT to docker-compose for image tagging.
        GIT_COMMIT=$(git rev-parse --short HEAD) \
            docker-compose -f factory/docker-compose.yml run --rm --name "ms-graph-worker-local" ms-graph-worker
        ;;

    "Trigger a remote job on Cloud Run")
        # For remote jobs, we only need to override the target entity.
        # The rest of the config is baked into the Cloud Run Job definition.
        local env_vars="TARGET_ENTITY=${TARGET_ENTITY}"

        gum style --bold "--> Triggering remote Cloud Run Job '${CLOUD_RUN_JOB_NAME}' to export '${TARGET_ENTITY}'..."
        gum spin --spinner dot --title "Executing remote job..." -- \
            gcloud run jobs execute "${CLOUD_RUN_JOB_NAME}" \
                --region "${GCP_REGION}" \
                --project "${GOOGLE_CLOUD_PROJECT}" \
                --update-env-vars="${env_vars}" \
                --wait
        ;;

    *)
        echo "Aborted."
        exit 0
        ;;
esac

echo
gum style --bold "✅ Run command complete."
