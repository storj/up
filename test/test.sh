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

go install -C ../

export STORJUP_NO_HISTORY=true

storj-up init minimal,db,uplink

docker compose down -v
docker compose up -d

storj-up health -d 90

docker compose exec -T uplink bash <<-'EOF'

  dd if=/dev/random of=data count=10240 bs=1024
  sha256sum data > sha256.sum

  eval $(storj-up credentials -s satellite-api:7777 -c satellite-api:10000 -e)

  BUCKET=bucket$RANDOM
  uplink --interactive=false mb sj://$BUCKET
  uplink --interactive=false cp data sj://$BUCKET/data

  rm data
  uplink --interactive=false cp sj://$BUCKET/data data
  sha256sum -c sha256.sum

EOF
