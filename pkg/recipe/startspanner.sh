#!/usr/bin/env bash
set -xem
echo "Starting Spanner dev container"

# Set default environment variables if not provided
export PROJECT_ID="${PROJECT_ID:-test-project}"
export INSTANCE_NAME="${INSTANCE_NAME:-test-instance}"

SPANNER_HTTP_URL=http://localhost:9020
MAX_RETRIES=5

# Start Spanner emulator with public access flags
gateway_main --hostname 0.0.0.0 --grpc_port 9010 --http_port 9020 &

spanner_post() {
  local url="$1"
  local data="$2"
  curl --silent --show-error --fail \
    --max-time 10 \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${data}" \
    "${url}"
}

create_instance() {
  for i in $(seq 1 $MAX_RETRIES); do
    if spanner_post \
      "${SPANNER_HTTP_URL}/v1/projects/${PROJECT_ID}/instances" \
      '{
        "instanceId": "'"${INSTANCE_NAME}"'",
        "instance": {
          "config": "projects/'"${PROJECT_ID}"'/instanceConfigs/emulator-config",
          "displayName": "Emulator with public access",
          "nodeCount": 1
        }
      }'
    then
      echo "Successfully created instance ${INSTANCE_NAME}"
      return 0
    fi
    echo "Failed to create instance, retrying in 5s... (attempt $i of ${MAX_RETRIES})"
    sleep 5
  done
  return 1
}

create_database() {
  local DB="$1"
  for i in $(seq 1 $MAX_RETRIES); do
    if spanner_post \
      "${SPANNER_HTTP_URL}/v1/projects/${PROJECT_ID}/instances/${INSTANCE_NAME}/databases" \
      '{
        "createStatement": "CREATE DATABASE `'"${DB}"'`"
      }'
    then
      echo "Successfully created database ${DB}"
      return 0
    fi
    echo "Failed to create database ${DB}, retrying in 5s... (attempt $i of ${MAX_RETRIES})"
    sleep 5
  done
  return 1
}

echo "Creating Spanner instance with public access"
if ! create_instance; then
  echo "Failed to create instance ${INSTANCE_NAME} after ${MAX_RETRIES} attempts. Exiting."
  exit 1
fi

echo "Creating Spanner databases with public access"
for DB in master metainfo satellite; do
  if ! create_database "${DB}"; then
    echo "Failed to create database ${DB} after ${MAX_RETRIES} attempts. Exiting."
    exit 1
  fi
done

# Make sure the environment variable is set for external access
echo "====================================================="
echo "Spanner emulator is running with the following config:"
echo "Project: ${PROJECT_ID}"
echo "Instance: ${INSTANCE_NAME}"
echo "HTTP Endpoint: ${SPANNER_HTTP_URL}"
echo "====================================================="

# Export for all processes in the container
export SPANNER_EMULATOR_HOST=0.0.0.0:9010

# Keep foreground process running
fg %1
