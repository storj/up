#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

cleanup() {
  if [ -f "docker-compose.yaml" ]
  then
    docker compose down
  fi
  rm -rf pass pk.json TestToken.abi TestToken.bin .contracts.yaml blockchain
  rm -rf docker-compose.yaml
}

trap cleanup EXIT

go install storj.io/storj-up

if [ ! "$(which storjscan )" ]; then
   go install storj.io/storjscan/cmd/storjscan@latest
fi

if [ ! "$(which cethacea)" ]; then
   go install github.com/elek/cethacea@latest
fi

export STORJUP_NO_HISTORY=true

storj-up init minimal,satellite-core,edge,db,billing
storj-up env setenv satellite-core STORJ_PAYMENTS_BILLING_CONFIG_INTERVAL=5s
storj-up env setenv satellite-core STORJ_PAYMENTS_STORJSCAN_INTERVAL=5s
storj-up env setenv satellite-core STORJ_PAYMENTS_STORJSCAN_CONFIRMATIONS=12
storj-up env setenv storjscan STORJ_TOKEN_PRICE_USE_TEST_PRICES=true

docker compose down -v
docker compose up -d

storj-up health -d 60

eval $(storj-up credentials -e)
COOKIE=$(storj-up credentials | grep -o 'Cookie.*')

export CETH_CHAIN=http://localhost:8545
export CETH_ACCOUNT=2e9a0761ce9815b95b2389634f6af66abe5fec2b1e04b772728442b4c35ea365
export CETH_CONTRACT=$(cethacea contract deploy --quiet --name TOKEN TestToken.bin --abi TestToken.abi '(uint256)' 1000000000000)

curl -X GET -u "eu1:eu1secret" http://127.0.0.1:12000/api/v0/auth/whoami
curl -X GET -u "us1:us1secret" http://127.0.0.1:12000/api/v0/auth/whoami

storjscan mnemonic >.mnemonic
storjscan generate >.wallets
storjscan import --input-file .wallets --api-key us1 --api-secret us1secret --address http://127.0.0.1:12000
storjscan mnemonic >.mnemonic
storjscan generate >.wallets
storjscan import --input-file .wallets --api-key eu1 --api-secret eu1secret --address http://127.0.0.1:12000
rm -rf .mnemonic .wallets

curl -X POST 'http://localhost:10000/api/v0/payments/wallet' --header "$COOKIE"
ADDRESS=$(curl -X GET -s http://localhost:10000/api/v0/payments/wallet --header "$COOKIE" | jq -r '.address')

#ACCOUNT is defined with environment variables above
for i in {1..15}; do cethacea token transfer 1000 "$ADDRESS"; done

storj-up health -t billing_transactions -n 3 -d 12

RESPONSE=$(curl -X GET http://localhost:10000/api/v0/payments/wallet/payments --header "$COOKIE")
STATUS=$(echo $RESPONSE | jq -r '.payments[-1].Status')

if [ "${STATUS}" != 'confirmed' ]; then
  echo "Test FAILED. Payment status: "${STATUS}""
  exit 1
fi
