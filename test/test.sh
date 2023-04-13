#!/usr/bin/env bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

cleanup() {
  if [ -f "docker-compose.yaml" ]
  then
    docker compose down
  fi
  rm -rf data sha256.sum
  rm -rf docker-compose.yaml
}

trap cleanup EXIT

go install storj.io/storj-up@latest

export STORJUP_NO_HISTORY=true

storj-up init minimal,db

docker compose down -v
docker compose up -d

storj-up health -d 60
dd if=/dev/random of=data count=10240 bs=1024
sha256sum data > sha256.sum

eval $(storj-up credentials -e)

BUCKET=bucket$RANDOM
uplink mb sj://$BUCKET
uplink cp data sj://$BUCKET/data

rm data
uplink cp sj://$BUCKET/data data 
sha256sum -c sha256.sum
