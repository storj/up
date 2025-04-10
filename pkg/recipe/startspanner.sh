#!/usr/bin/env bash
set -xem
echo "Starting Spanner dev container"

# Set default environment variables if not provided
export PROJECT_ID="${PROJECT_ID:-test-project}"
export INSTANCE_NAME="${INSTANCE_NAME:-test-instance}"
export SPANNER_EMULATOR_URL="${SPANNER_EMULATOR_URL:-http://localhost:9010}"

# Start Spanner emulator with public access flags
gcloud emulators spanner start --host-port=0.0.0.0:9010&

# Configure gcloud
gcloud config set disable_prompts true
gcloud config configurations create emulator 2>/dev/null || true
gcloud config set auth/disable_credentials true
gcloud config set project "${PROJECT_ID}"
gcloud config set api_endpoint_overrides/spanner "${SPANNER_EMULATOR_URL}"

# Set environment variable for emulator
export SPANNER_EMULATOR_HOST=localhost:9010

# Create instance with retries
echo "Creating Spanner instance with public access"
for i in {1..5}; do
  if gcloud spanner instances create "${INSTANCE_NAME}" --config=emulator-config --description="Emulator with public access" --nodes=1; then
    echo "Successfully created instance ${INSTANCE_NAME}"
    break
  else
    echo "Failed to create instance, retrying in 5s... (attempt $i of 5)"
    sleep 5
  fi
done

# Create databases with retries
echo "Creating Spanner databases with public access"
for DB in master metainfo satellite; do
  for i in {1..5}; do
    if gcloud spanner databases create "${DB}" --instance="${INSTANCE_NAME}"; then
      echo "Successfully created database ${DB}"
      break
    else
      echo "Failed to create database ${DB}, retrying in 5s... (attempt $i of 5)"
      sleep 5
    fi
  done
done

# Make sure the environment variable is set for external access
echo "====================================================="
echo "Spanner emulator is running with the following config:"
echo "Project: ${PROJECT_ID}"
echo "Instance: ${INSTANCE_NAME}"
echo "Endpoint: ${SPANNER_EMULATOR_URL}"
echo "====================================================="

# Export for all processes in the container
export SPANNER_EMULATOR_HOST=0.0.0.0:9010

# Keep foreground process running
fg %1
