#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

go install storj.io/storj-up

if [ ! "$(which uplink)" ]; then
   go install storj.io/storj/cmd/uplink@latest
fi

if [ ! "$(which rclone)" ]; then
  go install github.com/rclone/rclone@v1.56.2
fi

storj-up init

docker compose down -v
docker compose up -d

storj-up health

dd if=/dev/random of=data count=10240 bs=1024
sha256sum data > sha256.sum

eval $(storj-up credentials -e)
rclone config create --non-interactive storjdev3 storj access_grant=$UPLINK_ACCESS

# using internal satellite-api address
eval $(docker-compose exec satellite-api s3 storj-up credentials -e)
rclone config create --non-interactive storjdev3 s3 type=storj access_key_id=$AWS_ACCESS_KEY_ID secret_access_key=$AWS_SECRET_ACCESS_KEY endpoint=$STORJ_GATEWAY

BUCKET=bucket$RANDOM
rclone mkdir storjdevs3:$BUCKET
rclone copy data storjdevs3:$BUCKET/data
sha256sum -c sha256.sum

rm data
rclone copy storjdevs3:$BUCKET/data download
mv download/data ./
sha256sum -c sha256.sum
docker compose down -v
