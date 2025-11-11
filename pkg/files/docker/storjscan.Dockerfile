# syntax=docker/dockerfile:1.3
ARG TYPE
ARG SOURCE
FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/build:20251110-1 AS base

FROM base AS commit
ARG BRANCH
ARG COMMIT
RUN git clone https://github.com/storj/storjscan.git --branch ${BRANCH}
RUN cd storjscan && git reset --hard ${COMMIT}
WORKDIR storjscan

FROM base AS branch
ARG BRANCH
RUN git clone https://github.com/storj/storjscan.git --depth=1 --branch ${BRANCH}
WORKDIR storjscan

FROM ${SOURCE} AS github

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/storjscan.git
WORKDIR storjscan
RUN git fetch https://review.dev.storj.tools/storj/storjscan ${REF} && git checkout FETCH_HEAD

FROM base AS local
ARG PATH
WORKDIR /var/lib/storj/storjscan
COPY --chown=storj ${PATH} .

FROM --platform=$TARGETPLATFORM ${TYPE} AS binaries
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    go install -race ./cmd/...
RUN go install -race github.com/elek/cethacea@main

FROM img.dev.storj.io/storjup/base:20251110-1 AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
