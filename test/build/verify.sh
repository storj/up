#!/usr/bin/env bash
set -ueo pipefail

# Verify that storj-up init + build generates a valid docker-compose.yaml
# with proper project name and build args.

TMP=$(mktemp -d)
cd "$TMP"

echo "=== storj-up init ==="
/go/storj-up/storj-up init minimal,db

# Verify the generated compose file has a name and no version field.
grep -q '^name:' docker-compose.yaml || { echo "FAIL: docker-compose.yaml missing 'name:' field"; exit 1; }
if grep -q '^version:' docker-compose.yaml; then echo "FAIL: docker-compose.yaml has deprecated 'version:' field"; exit 1; fi

echo "=== storj-up build (satellite-api) ==="
/go/storj-up/storj-up build -s remote github satellite-api -c abc123

echo "=== storj-up build (storagenode) ==="
/go/storj-up/storj-up build -s remote github storagenode -c abc123

# Verify BUILD_TAG and BASE_TAG are set in the generated compose file.
grep -q 'BUILD_TAG:' docker-compose.yaml || { echo "FAIL: docker-compose.yaml missing BUILD_TAG"; exit 1; }
grep -q 'BASE_TAG:' docker-compose.yaml || { echo "FAIL: docker-compose.yaml missing BASE_TAG"; exit 1; }

# Verify the compose file is valid.
docker compose config -q

echo "=== docker compose config ==="
docker compose config

echo "PASS: storj-up build generates valid compose configuration"
