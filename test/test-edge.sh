#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

if [ ! "$(which storj-up)" ]; then
   go install storj.io/storj-up@latest
fi

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

storj-up credentials -w

BUCKET=bucket$RANDOM
rclone mkdir storjdevs3:$BUCKET
rclone copy data storjdevs3:$BUCKET/data
sha256sum -c sha256.sum

rm data
rclone copy storjdevs3:$BUCKET/data download
mv download/data ./
sha256sum -c sha256.sum
docker compose down -v
