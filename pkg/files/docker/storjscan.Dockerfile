# syntax=docker/dockerfile:1.3
ARG TYPE
ARG SOURCE
ARG BUILD_TAG
ARG BASE_TAG
FROM --platform=$TARGETPLATFORM ghcr.io/storj/up/build:${BUILD_TAG} AS base

FROM base AS commit
ARG BRANCH
ARG COMMIT
RUN git clone https://github.com/storj/storjscan.git --branch ${BRANCH}
RUN cd storjscan && git reset --hard ${COMMIT}
WORKDIR /storj/storjscan

FROM base AS branch
ARG BRANCH
RUN git clone https://github.com/storj/storjscan.git --depth=1 --branch ${BRANCH}
WORKDIR /storj/storjscan

FROM ${SOURCE} AS github

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/storjscan.git
WORKDIR /storj/storjscan
RUN git fetch https://review.dev.storj.tools/storj/storjscan ${REF} && git checkout FETCH_HEAD

FROM base AS local
ARG PATH
WORKDIR /storj/storjscan
COPY --chown=storj ${PATH} .

FROM --platform=$TARGETPLATFORM ${TYPE} AS binaries
RUN --mount=type=cache,target=/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/storj/.cache/go-build,mode=777,uid=1000 \
    go install -race ./cmd/...
RUN go install -race github.com/elek/cethacea@main

ARG BASE_TAG
FROM ghcr.io/storj/up/base:${BASE_TAG} AS final
COPY --from=binaries /storj/bin /storj/bin
LABEL org.opencontainers.image.source="https://github.com/storj/up"
