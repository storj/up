#!/usr/bin/env bash
set -xem
echo "Starting Spanner dev container"
gcloud emulators spanner start --host-port=0.0.0.0:9010&
gcloud config set disable_prompts true
gcloud config configurations create emulator
gcloud config set auth/disable_credentials true
gcloud config set project "${PROJECT_ID}"
gcloud config set api_endpoint_overrides/spanner "${SPANNER_EMULATOR_URL}"
gcloud spanner instances create "${INSTANCE_NAME}" --config=emulator-config --description=Emulator --nodes=1

#TODO: find a more flexible way to create the databases
gcloud spanner databases create master --instance="${INSTANCE_NAME}"
gcloud spanner databases create metainfo --instance="${INSTANCE_NAME}"
gcloud spanner databases create satellite --instance="${INSTANCE_NAME}"

fg %1
