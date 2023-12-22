dd if=/dev/random of=data count=10240 bs=1024
sha256sum data > sha256.sum

curl https://rclone.org/install.sh | bash
# shellcheck disable=SC2086
rclone config create storjdevs3 s3 env_auth true provider Minio access_key_id $AWS_ACCESS_KEY_ID secret_access_key "$AWS_SECRET_ACCESS_KEY" endpoint "$STORJ_GATEWAY" chunk_size 64M upload_cutoff 64M use_already_exists false

BUCKET=bucket$RANDOM
rclone mkdir storjdevs3:$BUCKET
rclone copy data storjdevs3:$BUCKET/data
sha256sum -c sha256.sum

rm data
rclone copy storjdevs3:$BUCKET/data download
mv download/data ./
sha256sum -c sha256.sum