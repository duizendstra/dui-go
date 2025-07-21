#!/bin/bash
# factory/scripts/release.sh
# This script performs a versioned release, including safety checks,
# git tagging, building, pushing, and deploying.

set -e

# --- 1. Input Validation ---
VERSION=$1
if [ -z "$VERSION" ]; then
  gum style --bold --foreground 196 -- "Error: No version tag provided." >&2
  echo "Usage: task release -- v1.2.3" >&2
  exit 1
fi

gum style --border double --border-foreground 212 --padding "1 2" -- "--- 📦 Releasing Version: $(gum style --bold "$VERSION") ---"

# --- 2. Pre-flight Safety Checks ---
gum style --bold --padding "0 1" -- "--- 1. Running Pre-flight Safety Checks ---"

if ! git diff --quiet --exit-code; then
  gum style --bold --foreground 196 -- "Error: You have uncommitted changes. Please commit or stash them before releasing." >&2
  exit 1
fi
echo "✅ Git working directory is clean."

if git rev-parse -q --verify "refs/tags/${VERSION}" >/dev/null; then
  gum style --bold --foreground 196 -- "Error: Git tag '${VERSION}' already exists." >&2
  exit 1
fi
echo "✅ Git tag '${VERSION}' is available."

gum spin --spinner dot --title "Running unit tests..." -- \
    task test -- --unit
echo "✅ All unit tests passed."

gum style --bold -- "✅ Pre-flight checks complete."

# --- 3. Tagging and Deployment ---
gum style --bold --padding "0 1" -- "--- 2. Building, Pushing, and Deploying ---"
# These variables are constructed from the .env file to ensure consistency.
readonly ARTIFACT_REGISTRY_REPO="${GCP_REGION}-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/${ARTIFACT_REGISTRY_NAME}"
readonly IMAGE_FQN="${ARTIFACT_REGISTRY_REPO}/${CONTAINER_NAME}"
# For a release, Workload Identity is mandatory.
readonly ENV_VARS_STRING="^;^GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT};GCS_BUCKET_NAME=${GCS_BUCKET_NAME};SECRET_NAME=${SECRET_NAME};LOG_LEVEL=${LOG_LEVEL};USE_WORKLOAD_IDENTITY=true;SERVICE_NAME=${SERVICE_NAME}"

echo "--> Creating git tag '${VERSION}'..."
git tag -a "${VERSION}" -m "Release ${VERSION}"

echo "--> Building and pushing container image: ${IMAGE_FQN}:${VERSION}"
gum spin --spinner dot --title "Building container..." -- \
docker buildx build \
  --file factory/Dockerfile \
  --platform linux/amd64 \
  -t "${IMAGE_FQN}:${VERSION}" \
  -t "${IMAGE_FQN}:latest" \
  --push \
  .

echo "--> Deploying job '${CLOUD_RUN_JOB_NAME}' with image tag '${VERSION}'..."
gum spin --spinner dot --title "Deploying to Cloud Run..." -- \
gcloud run jobs deploy "${CLOUD_RUN_JOB_NAME}" \
  --image "${IMAGE_FQN}:${VERSION}" \
  --region "${GCP_REGION}" \
  --service-account "${CLOUD_RUN_SA}" \
  --set-env-vars="${ENV_VARS_STRING}" \
  --project="${GOOGLE_CLOUD_PROJECT}" \
  --quiet

# --- 4. Finalization ---
gum style --bold --padding "0 1" -- "--- 3. Finalizing Release ---"
echo "--> Pushing git tag to remote..."
gum spin --spinner dot --title "Pushing git tag..." -- git push origin "${VERSION}"

echo
gum style --bold -- "✅ Release ${VERSION} is complete."
