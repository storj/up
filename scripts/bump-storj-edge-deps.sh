#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_DIR
ROOT_DIR=$(cd "${SCRIPT_DIR}"/.. && pwd)
readonly ROOT_DIR
COFIG_GEN_DIR="${ROOT_DIR}/pkg/config/gen"
readonly COFIG_GEN_DIR

fail() {
  msg="${1}"
  code="${2:-1}"
  echo "${msg}"
  exit "${code}"
}

assert_version_prefix() {
  version="${1}"
  if [ "${version#v}" == "${version}" ]; then
    fail "a version must start with v"
  fi
}

if [ "$#" -gt 2 ]; then
  echo "invalid number of arguments. Use: ${0} <storj-version> [<edge-version>]"
fi

if [ "$#" -lt 1 ]; then
  fail "storj.io/storj version is required"
fi

storj_version="${1}"
assert_version_prefix "${storj_version}"

edge_version=""
if [ "$#" -eq 2 ]; then
  edge_version="${2}"
  assert_version_prefix "${edge_version}"
fi

pushd "${ROOT_DIR}"
go get storj.io/storj@"${storj_version}"
go mod tidy
popd

pushd "${COFIG_GEN_DIR}"
go get storj.io/storj@"${storj_version}"

if [ -n "${edge_version}" ]; then
  go get storj.io/edge@"${edge_version}"
fi
go mod tidy
popd

go generate "${COFIG_GEN_DIR}/.."

if ! go run . -h 2>&1 > /dev/null; then
  fail "all the automated changes are applied, but the version bump required manual changes, please apply them and verify that the tool compiles" 2
else
  echo "Version bump succeeded"
fi
