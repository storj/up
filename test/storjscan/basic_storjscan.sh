set -xueo pipefail

export CETH_CHAIN=http://geth:8545
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

curl -X POST 'http://satellite-api:10000/api/v0/payments/wallet' --header "Cookie: _tokenKey=$_tokenKey"
ADDRESS=$(curl -X GET -s http://satellite-api:10000/api/v0/payments/wallet --header "Cookie: _tokenKey=$_tokenKey" | jq -r '.address')

#ACCOUNT is defined with environment variables above
for i in {1..15}; do cethacea token transfer 1000 "$ADDRESS"; sleep 1; done
storj-up health -t billing_transactions -n 6 -d 12

RESPONSE=$(curl -X GET http://satellite-api:10000/api/v0/payments/wallet/payments --header "Cookie: _tokenKey=$_tokenKey")
STATUS=$(echo "$RESPONSE" | jq -r '.payments[-4].Status')
STATUS_BONUS=$(echo "$RESPONSE" | jq -r '.payments[-1].Status')

if [ "${STATUS_BONUS}" != 'complete' ] || [ "${STATUS}" != 'confirmed' ]; then
  echo "Test FAILED. Payment status: ${STATUS} Payment bonus status: ${STATUS_BONUS}"
  exit 1
else
  echo "Test PASSED."
fi
