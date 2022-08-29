#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

if [ ! "$(which storj-up)" ]; then
   go install storj.io/storj-up@latest
fi

if [ ! "$(which storjscan )" ]; then
   go install storj.io/storjscan/cmd/storjscan@latest
fi

if [ ! "$(which cethacea)" ]; then
   go install github.com/elek/cethacea@latest
fi

export STORJUP_NO_HISTORY=true

storj-up init storj,db,billing

docker compose down -v
docker compose up -d

storj-up health

eval $(storj-up credentials -e)
COOKIE=$(storj-up credentials | grep -o 'Cookie.*')

export CETH_CHAIN=http://localhost:8545
export CETH_ACCOUNT=2e9a0761ce9815b95b2389634f6af66abe5fec2b1e04b772728442b4c35ea365
export CETH_CONTRACT=$(cethacea contract deploy --quiet --name TOKEN test-blockchain/TestToken.bin --abi test-blockchain/TestToken.abi '(uint256)' 1000000000000)

curl -X GET -u "eu1:eu1secret" http://127.0.0.1:12000/api/v0/auth/whoami
curl -X GET -u "us1:us1secret" http://127.0.0.1:12000/api/v0/auth/whoami

storjscan mnemonic > .mnemonic
storjscan generate --api-key us1 --api-secret us1secret --address http://127.0.0.1:12000
storjscan mnemonic > .mnemonic
storjscan generate --api-key eu1 --api-secret eu1secret --address http://127.0.0.1:12000
rm -rf .mnemonic

curl --location --request POST 'http://localhost:10000/api/v0/payments/wallet' --header "$COOKIE"
ADDRESS=$(curl -sb GET http://localhost:10000/api/v0/payments/wallet --header "$COOKIE" | jq -r '.address')

#ACCOUNT is defined with environment variables above
for i in {1..15}; do cethacea token transfer 1000 0x"$ADDRESS"; done

# todo: find a better way than sleep to wait until token balance chores reflect in billing table
sleep 180

curl -sb GET http://localhost:10000/api/v0/payments/wallet --header "$COOKIE"

rm -rf test-blockchain
docker compose down
rm -rf docker-compose.yaml
