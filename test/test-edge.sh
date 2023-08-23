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

go install -C ../

export STORJUP_NO_HISTORY=true

storj-up init minimal,edge,db,uplink

docker compose down -v
docker compose up -d

storj-up health -d 90

docker compose exec -T -u 0 uplink bash <<-'EOF'

  dd if=/dev/random of=data count=10240 bs=1024
  sha256sum data > sha256.sum

  eval $(storj-up credentials -s satellite-api:7777 -c satellite-api:10000 -a http://authservice:8888 -e --s3)
  #todo: add curl and unzip to base image
  apt-get update
  apt-get -y install curl unzip
  curl https://rclone.org/install.sh | bash
  rclone config create storjdevs3 s3 env_auth true provider Minio access_key_id $AWS_ACCESS_KEY_ID secret_access_key $AWS_SECRET_ACCESS_KEY endpoint http://gateway-mt:9999 chunk_size 64M upload_cutoff 64M

  BUCKET=bucket$RANDOM
  rclone mkdir storjdevs3:$BUCKET
  rclone copy data storjdevs3:$BUCKET/data
  sha256sum -c sha256.sum

  rm data
  rclone copy storjdevs3:$BUCKET/data download
  mv download/data ./
  sha256sum -c sha256.sum

EOF
