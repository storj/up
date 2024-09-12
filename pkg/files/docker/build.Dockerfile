# syntax=docker/dockerfile:1.3
FROM --platform=$TARGETPLATFORM ubuntu:22.04 as base
ARG TARGETPLATFORM
RUN apt-get update
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install curl && curl -sfL https://deb.nodesource.com/setup_22.x  | bash -
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install git sudo nodejs make gcc brotli g++
RUN echo ${TARGETPLATFORM} | sed 's/linux\///' | xargs -I PLATFORM curl --fail -L https://go.dev/dl/go1.22.7.linux-PLATFORM.tar.gz | tar -C /usr/local -xz && cp /usr/local/go/bin/go /usr/local/bin/go
ENV GOROOT=/usr/local/go

RUN useradd storj --uid 1000 -d /var/lib/storj && \
    mkdir -p /var/lib/storj/shared && \
    mkdir -p /var/lib/storj/.cache && \
    mkdir -p /var/lib/storj/go/pkg/mod && \
    chown -R storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000  \
    go install github.com/go-delve/delve/cmd/dlv@latest
ADD pkg/recipe/entrypoint.sh /var/lib/storj/entrypoint.sh

