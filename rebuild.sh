#!/usr/bin/env bash
set -ex
docker build -t ghcr.io/elek/storj -f storj.Dockerfile .
docker build -t ghcr.io/elek/storj-edge -f edge.Dockerfile .
