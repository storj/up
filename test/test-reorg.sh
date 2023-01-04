#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex
set -o pipefail  # make a failure on the left-hand side fail the entire pipeline

wait_for_geth() {
while ! { curl -s -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67}' http://localhost:8545 \
          | jq -e '.result' >/dev/null; }; do
    echo "Waiting for Geth to start. Trying again in 2 seconds."
    sleep 2
done
}

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

storj-up init minimal,satellite-core,satellite-admin,edge,db,billing
storj-up env setenv satellite-core STORJ_PAYMENTS_BILLING_CONFIG_INTERVAL=5s
storj-up env setenv satellite-core STORJ_PAYMENTS_STORJSCAN_INTERVAL=5s
storj-up env setenv satellite-core STORJ_PAYMENTS_STORJSCAN_CONFIRMATIONS=12
storj-up env setenv storjscan STORJ_TOKEN_PRICE_USE_TEST_PRICES=true

docker compose down -v
docker compose up -d

storj-up health

eval $(storj-up credentials -e)
COOKIE=$(storj-up credentials | grep -o 'Cookie.*')

export CETH_CHAIN=http://localhost:8545
export CETH_ACCOUNT=2e9a0761ce9815b95b2389634f6af66abe5fec2b1e04b772728442b4c35ea365
export CETH_CONTRACT=$(cethacea contract deploy --quiet --name TOKEN TestToken.bin --abi TestToken.abi '(uint256)' 1000000000000)

storjscan mnemonic > .mnemonic
storjscan generate >.wallets
storjscan import --input-file .wallets --api-key us1 --api-secret us1secret --address http://127.0.0.1:12000
rm -rf .mnemonic .wallets

curl -X POST 'http://localhost:10000/api/v0/payments/wallet' --header "$COOKIE"
ADDRESS=$(curl -X GET -s http://localhost:10000/api/v0/payments/wallet --header "$COOKIE" | jq -r '.address')

#15 transactions means 3 are fully confirmed
for i in {1..15}; do cethacea token transfer 1000 "$ADDRESS"; done
storj-up health -t billing_transactions -n 3 -d 12

#save the last transaction of the base chain
CONTENT=$(curl -s -X GET http://localhost:10000/api/v0/payments/wallet/payments --header "$COOKIE")
BASE0=$(jq -r '.payments[0].ID' <<< "${CONTENT}")

#save off the base chain and restart geth
docker compose stop geth
cp -rf blockchain/geth blockchain/base-chain
docker compose up -d

wait_for_geth

#adding 5 transactions for a total of 20 means 8 should be fully confirmed
for i in {1..5}; do cethacea token transfer 1 "$ADDRESS"; done
storj-up health -t billing_transactions -n 8 -d 12

#save the last transaction (0) and compare to the base chain
CONTENT=$(curl -s -X GET http://localhost:10000/api/v0/payments/wallet/payments --header "$COOKIE")
PREREORG4=$(jq -r '.payments[4].ID' <<< "${CONTENT}")
PREREORG5=$(jq -r '.payments[5].ID' <<< "${CONTENT}")
if [[ ${BASE0%??} -ne ${PREREORG5%??} ]]; then
  exit 1
fi

#resetting back to the base chain, so the extra 5 transactions above disappear as a "reorg" would cause
docker compose stop geth
rm -rf blockchain/geth
mv blockchain/base-chain blockchain/geth
docker compose up -d

wait_for_geth

storj-up health -t storjscan_payments -n 15 -d 12

#adding 5 different transactions to simulate reorg
for i in {1..5}; do cethacea token transfer 1000 "$ADDRESS"; done
storj-up health -t storjscan_payments -n 20 -d 12

CONTENT=$(curl -s -X GET http://localhost:10000/api/v0/payments/wallet/payments --header "$COOKIE")
POSTREORG4=$(jq -r '.payments[4].ID' <<< "${CONTENT}")
POSTREORG5=$(jq -r '.payments[5].ID' <<< "${CONTENT}")
if [[ ${BASE0%??} -ne ${POSTREORG5%??} ]]; then
  exit 1
fi

# should be different IDs for transactions after the base
if [[ ${PREREORG4%??} == ${POSTREORG4%??} ]]; then
  exit 1
fi
