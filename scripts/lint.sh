#!/usr/bin/env bash
RESULT=0

function runtest(){
echo "=======$@======="
if ! "$@"; then
   echo "FAILED"
   RESULT=1
fi
}

runtest check-copyright
runtest check-large-files
runtest check-imports -race ./...
runtest check-peer-constraints -race
runtest check-atomic-align ./...
runtest check-errs ./...
runtest check-monkit ./...
runtest staticcheck ./...
runtest golangci-lint --build-tags mage -j=2 run
runtest check-mod-tidy -mod .build/go.mod.orig
exit $RESULT
