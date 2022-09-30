#!/usr/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

set -ex

cleanup() {
  if [ -f "docker-compose.yaml" ]
  then
    docker compose down
  fi
  rm -rf .contracts.yaml
  rm -rf storjscan
  rm -rf geth
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

storj-up init storj,db,billing
storj-up env setenv satellite-api satellite-core satellite-admin STORJ_PAYMENTS_PROVIDER=stripecoinpayments
storj-up env setenv satellite-api satellite-core satellite-admin STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_STRIPE_PUBLIC_KEY="$STRIPE_PUBLIC_KEY"
storj-up env setenv satellite-api satellite-core satellite-admin STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_STRIPE_SECRET_KEY="$STRIPE_SECRET_KEY"
storj-up env setenv satellite-api satellite-core satellite-admin STORJ_PAYMENTS_BILLING_CONFIG_INTERVAL=5s
storj-up env setenv satellite-api satellite-core satellite-admin STORJ_PAYMENTS_STORJSCAN_INTERVAL=5s
storj-up env setenv satellite-api satellite-core satellite-admin STORJ_PAYMENTS_STORJSCAN_CONFIRMATIONS=12

docker compose up -d

storj-up health

eval $(storj-up credentials -e)
COOKIE=$(storj-up credentials | grep -o 'Cookie.*')

export CETH_CHAIN=http://localhost:8545
export CETH_ACCOUNT=2e9a0761ce9815b95b2389634f6af66abe5fec2b1e04b772728442b4c35ea365
export CETH_CONTRACT=$(cethacea contract deploy --quiet --name TOKEN storjscan/test-contract/TestToken.bin --abi storjscan/test-contract/TestToken.abi '(uint256)' 1000000000000)

curl -X GET -u "eu1:eu1secret" http://127.0.0.1:12000/api/v0/auth/whoami
curl -X GET -u "us1:us1secret" http://127.0.0.1:12000/api/v0/auth/whoami

storjscan mnemonic > .mnemonic
storjscan generate --api-key us1 --api-secret us1secret --address http://127.0.0.1:12000
storjscan mnemonic > .mnemonic
storjscan generate --api-key eu1 --api-secret eu1secret --address http://127.0.0.1:12000
rm -rf .mnemonic

curl -X POST 'http://localhost:10000/api/v0/payments/wallet' --header "$COOKIE"
ADDRESS=$(curl -X GET -s http://localhost:10000/api/v0/payments/wallet --header "$COOKIE" | jq -r '.address')

#ACCOUNT is defined with environment variables above
for i in {1..15}; do cethacea token transfer 1000000000 0x"$ADDRESS"; done

storj-up health -t billing_transactions -n 3 -d 12

curl -X GET http://localhost:10000/api/v0/payments/wallet --header "$COOKIE"
curl -X POST http://localhost:10000/api/v0/payments/account --header "$COOKIE"
STRIPE_CUSTOMERS=$(docker exec test-cockroach-1 cockroach sql --insecure -d master -e "SELECT customer_id FROM stripe_customers";)
CUSTOMER_ID=$(echo "$STRIPE_CUSTOMERS" | grep -o 'cus_[^[:blank:]]*')
PAYMENT_INTENT=$(curl https://api.stripe.com/v1/payment_intents -u "${STRIPE_SECRET_KEY}:" -d customer="$CUSTOMER_ID" -d amount=16742 -d currency=usd -d payment_method=pm_card_visa -d setup_future_usage=off_session )
PAYMENT_METHOD=$(jq -r '.payment_method' <<< "${PAYMENT_INTENT}")
curl https://api.stripe.com/v1/payment_methods/"$PAYMENT_METHOD"/attach -u "$STRIPE_SECRET_KEY": -d customer="$CUSTOMER_ID"
curl https://api.stripe.com/v1/customers/"$CUSTOMER_ID" -u "${STRIPE_SECRET_KEY}:" -d "invoice_settings[default_payment_method]"="$PAYMENT_METHOD"

# invoicing
storj-up testdata project-usage
DAY=$(date +%d)
MONTH=$(date +%m)
YEAR=$(date +%Y)
if [ $(uname -s) == "Darwin" ]
then
  LAST_MONTH=$(date -v-1m +%m)
  LAST_MONTH_YEAR=$(date -v-1m +%Y)
else
  LAST_MONTH=$(date -d "$(date +%Y-%m-1) -1 month" +%m)
  LAST_MONTH_YEAR=$(date -d "$(date +%Y-%m-1) -1 month" +%Y)
fi

docker compose exec satellite-admin satellite billing prepare-invoice-records "$LAST_MONTH"/"$LAST_MONTH_YEAR" --log.level=info --log.output=stdout
docker compose exec satellite-admin satellite billing create-project-invoice-items "$LAST_MONTH"/"$LAST_MONTH_YEAR" --log.level=info --log.output=stdout
docker compose exec satellite-admin satellite billing create-invoices "$LAST_MONTH"/"$LAST_MONTH_YEAR" --log.level=info --log.output=stdout
docker compose exec satellite-admin satellite billing finalize-invoices --log.level=info --log.output=stdout
docker compose exec satellite-admin satellite billing pay-invoices "$DAY"/"$MONTH"/"$YEAR" --log.level=info --log.output=stdout

BALANCE=$(curl -X GET -s http://localhost:10000/api/v0/payments/wallet --header "$COOKIE" | jq -r '.balance')

if [[ $BALANCE == -* ]]
then
  exit 1
fi
