#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

if [ ! "$(which sjr)" ]; then
   go install github.com/elek/sjr@latest
fi

if [ ! "$(which uplink)" ]; then
   go install storj.io/storj/cmd/uplink@latest
fi

sjr init minimal db
sjr scale 10 storagenode

docker-compose down -v
docker-compose up -d

sjr health
dd if=/dev/random of=data count=10240 bs=1024
sha256sum data > sha256.sum

eval $(sjr credentials -e)

BUCKET=bucket$RANDOM
uplink mb sj://$BUCKET
uplink cp data sj://$BUCKET/data

rm data
uplink cp sj://$BUCKET/data data 
sha256sum -c sha256.sum
docker-compose down
