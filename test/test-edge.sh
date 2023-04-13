#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

cleanup() {
  if [ -f "docker-compose.yaml" ]
  then
    docker compose down
  fi
  rm -rf data download sha256.sum
  rm -rf docker-compose.yaml
}

trap cleanup EXIT

go install storj.io/storj-up@latest

export STORJUP_NO_HISTORY=true

storj-up init

docker compose down -v
docker compose up -d

storj-up health -d 60

dd if=/dev/random of=data count=10240 bs=1024
sha256sum data > sha256.sum

eval $(storj-up credentials -e)
rclone config create --non-interactive storjdev3 storj access_grant=$UPLINK_ACCESS

# using internal satellite-api address
eval $(docker-compose exec -T satellite-api storj-up credentials --s3 -e -a http://authservice:8888 -s satellite-api:7777)
rclone config create --non-interactive storjdevs3 s3 type=s3 provider=Storj access_key_id=$AWS_ACCESS_KEY_ID secret_access_key=$AWS_SECRET_ACCESS_KEY endpoint=http://localhost:9999

BUCKET=bucket$RANDOM
rclone mkdir storjdevs3:$BUCKET
rclone copy data storjdevs3:$BUCKET/data
sha256sum -c sha256.sum

rm data
rclone copy storjdevs3:$BUCKET/data download
mv download/data ./
sha256sum -c sha256.sum
