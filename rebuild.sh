#!/usr/bin/env bash
set -ex
#docker build -t ghcr.io/elek/storj-build -f build.Dockerfile .
docker build -t ghcr.io/elek/storj:1.39.6 --build-arg BRANCH=v1.39.6 -f pkg/storj.Dockerfile .
docker build -t ghcr.io/elek/storj-edge:1.14.3 --build-arg BRANCH=v1.14.3  -f pkg/edge.Dockerfile .
