dd if=/dev/random of=data count=10240 bs=1024
sha256sum data > sha256.sum

BUCKET=buckett$RANDOM
uplink --interactive=false mb sj://$BUCKET
uplink --interactive=false cp data sj://$BUCKET/data

rm data
uplink --interactive=false cp sj://$BUCKET/data data
sha256sum -c sha256.sum