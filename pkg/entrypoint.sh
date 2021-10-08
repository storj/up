#!/usr/bin/env bash
set -xe
echo "Starting Storj development container"
#Identifying IP address
export STORJ_NODE_IP=$(ip route get 8.8.8.8 | awk -F"src " 'NR==1{split($2,a," ");print a[1]}')

#Generate identity if missing
if [ "$STORJ_IDENTITY_DIR" ]; then
  if [ ! -f "$STORJ_IDENTITY_DIR/identity.key" ]; then
    identity --identity-dir $STORJ_IDENTITY_DIR --difficulty 8 create .
  fi
fi

if [ "$STORJ_WAIT_FOR_DB" ]; then
  devrun wait-for-port cockroach:26257
  devrun wait-for-port redis:6379
fi

if [ "$STORJ_WAIT_FOR_SATELLITE" ]; then
  SATELLITE_ADDRESS=$(devrun wait-for-satellite satellite-api:7777)
fi


if [ "$STORJ_ROLE" == "satellite-api" ]; then
  mkdir -p /var/lib/storj/.local

  #only migrate first time
  if [ ! -f "/var/lib/storj/.local/migrated" ]; then
    satellite run migration --identity-dir $STORJ_IDENTITY_DIR
    touch /var/lib/storj/.local/migrated

  fi
fi

if [ "$STORJ_ROLE" == "storagenode" ]; then
  #Initialize config, required only to have all the dirs created
  export STORJ_CONTACT_EXTERNAL_ADDRESS=$STORJ_NODE_IP:28967
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
    devrun credentials satellite-api test@storj.io
    eval $(devrun credentials grant satellite-api test@storj.io)
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
