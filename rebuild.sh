#!/usr/bin/env bash
set -ex
docker build -t elek/storj -f storj.Dockerfile .
docker build -t elek/storj-edge -f edge.Dockerfile .
