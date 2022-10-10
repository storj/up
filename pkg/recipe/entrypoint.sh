#!/usr/bin/env bash
set -xe
echo "Starting Storj development container"
#Identifying IP address
export STORJ_NODE_IP=$(ip route get 8.8.8.8 | awk -F"src " 'NR==1{split($2,a," ");print a[1]}')

#Generate identity if missing
if [ "$STORJ_IDENTITY_DIR" ]; then
  if [ ! -f "$STORJ_IDENTITY_DIR/identity.key" ]; then
    if [ "$STORJ_USE_PREDEFINED_IDENTITY" ]; then
      # use predictable, pre-generated identity
      mkdir -p $(dirname $STORJ_IDENTITY_DIR)
      cp -r /var/lib/storj/identities/$STORJ_USE_PREDEFINED_IDENTITY $STORJ_IDENTITY_DIR
    else
      identity --identity-dir $STORJ_IDENTITY_DIR --difficulty 8 create .
    fi
  fi
fi

for i in ${STORJ_WAIT_FOR//,/ }
do
    storj-up util wait-for-port "$i"
done

if [ "$STORJ_WAIT_FOR_DB" ]; then
  storj-up util wait-for-port cockroach:26257
  storj-up util wait-for-port redis:6379
fi

if [ "$STORJ_WAIT_FOR_SATELLITE" ]; then
  SATELLITE_ADDRESS=$(storj-up util wait-for-satellite satellite-api:7777)
fi


if [ "$STORJ_ROLE" == "satellite-api" ]; then
  mkdir -p /var/lib/storj/.local

  #only migrate first time, if automatic migration is not defined
  if [[ ! -f "/var/lib/storj/.local/migrated" && -z "$STORJ_DATABASE_OPTIONS_MIGRATION_UNSAFE" ]]; then
    satellite run migration --identity-dir $STORJ_IDENTITY_DIR
    touch /var/lib/storj/.local/migrated
  fi
fi

if [ "$STORJ_ROLE" == "storagenode" ]; then
  #Initialize config, required only to have all the dirs created
  : ${STORJ_CONTACT_EXTERNAL_ADDRESS:=$STORJ_NODE_IP:28967}
  if [ -f "/var/lib/storj/.local/share/storj/storagenode/config.yaml" ]; then
    rm "/var/lib/storj/.local/share/storj/storagenode/config.yaml"
  fi
  storagenode setup
fi

if [ "$STORJ_ROLE" == "multinode" ]; then
  if [ -f "/var/lib/storj/.local/share/storj/multinode/config.yaml" ]; then
    rm "/var/lib/storj/.local/share/storj/multinode/config.yaml"
  fi
  multinode setup
fi

if [ "$STORJ_ROLE" == "uplink" ]; then
  if [ "$1" != "/usr/bin/sleep"  ]; then
    storj-up credentials -e satellite-api test@storj.io
    eval $(storj-up credentials -e satellite-api test@storj.io)
  fi
fi
mkdir -p /var/lib/storj/.local/share/storj/satellite

if [ "$GO_DLV" ]; then
  echo "Starting with go dlv"

  #absolute file path is required
  CMD=$(which $1)
  shift
  /var/lib/storj/go/bin/dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec --check-go-version=false -- $CMD "$@"
else
   exec "$@"
fi
