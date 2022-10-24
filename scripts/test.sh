#!/usr/bin/env bash
set -exuo pipefail
cd $(dirname "${BASH_SOURCE[0]}")/..
mkdir -p build
go test -json ./... | tee build/tests.json | jq '. | select(.Action == "fail")' -cr 

