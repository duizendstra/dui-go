#!/bin/bash
# factory/scripts/deploy.sh
# This script orchestrates the entire deployment process. It builds, pushes,
# and then deploys the container to Cloud Run.

# Exit immediately if any command fails.
set -e

gum style --border double --border-foreground 212 --padding "1 2" -- "--- 🚀 Application Deployment ---"

# --- 1. Define Variables ---
gum style --bold --padding "0 1" -- "--- 1. Gathering Deployment Variables ---"

# These variables are sourced from the .env file by the root Taskfile.
readonly ARTIFACT_REGISTRY_REPO="${GCP_REGION}-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/${ARTIFACT_REGISTRY_NAME}"
readonly GIT_COMMIT=$(git rev-parse --short HEAD)
readonly IMAGE_FQN="${ARTIFACT_REGISTRY_REPO}/${CONTAINER_NAME}"

# For deployments, Workload Identity is mandatory.
readonly ENV_VARS_STRING="^;^GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT};GCS_BUCKET_NAME=${GCS_BUCKET_NAME};SECRET_NAME=${SECRET_NAME};LOG_LEVEL=${LOG_LEVEL};USE_WORKLOAD_IDENTITY=true;SERVICE_NAME=${SERVICE_NAME}"

echo "--> Image will be tagged as: ${IMAGE_FQN}:${GIT_COMMIT}"
echo "--> Job will be deployed to project: ${GOOGLE_CLOUD_PROJECT}"
echo "--> Job will run as service account: ${CLOUD_RUN_SA}"

# --- 2. Build and Push Container ---
echo
gum style --bold --padding "0 1" -- "--- 2. Building & Pushing Container Image ---"
echo "--> Authenticating Docker with Google Artifact Registry..."
gcloud auth configure-docker "${GCP_REGION}-docker.pkg.dev" --quiet

echo "--> Building for linux/amd64 and pushing..."
docker buildx build \
  --file factory/Dockerfile \
  --platform linux/amd64 \
  -t "${IMAGE_FQN}:${GIT_COMMIT}" \
  -t "${IMAGE_FQN}:latest" \
  --push \
  .

# --- 3. Deploy to Cloud Run ---
echo
gum style --bold --padding "0 1" -- "--- 3. Deploying to Google Cloud Run ---"
echo "--> Deploying job '${CLOUD_RUN_JOB_NAME}' to region '${GCP_REGION}'..."
gum spin --spinner dot --title "Deploying to Cloud Run..." -- gcloud run jobs deploy "${CLOUD_RUN_JOB_NAME}" \
  --image "${IMAGE_FQN}:${GIT_COMMIT}" \
  --region "${GCP_REGION}" \
  --service-account "${CLOUD_RUN_SA}" \
  --set-env-vars="${ENV_VARS_STRING}" \
  --project="${GOOGLE_CLOUD_PROJECT}" \
  --quiet

echo
gum style --bold -- "✅ Deployment complete."
