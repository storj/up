#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

if [ ! "$(which sjr)" ]; then
   go install github.com/elek/sjr@latest
fi

if [ ! "$(which uplink)" ]; then
   go install storj.io/storj/cmd/uplink@latest
fi


sjr init

docker-compose down -v
docker-compose up -d --scale storagenode=10
#yes it's a big todo, let's check if all the storagenodes are registered
sleep 60
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
