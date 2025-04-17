# syntax=docker/dockerfile:1.3
ARG TYPE
ARG SOURCE
FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/build:20240911-1 AS base

FROM base AS commit
ARG BRANCH
ARG COMMIT
RUN git clone https://github.com/storj/gateway-mt.git --branch ${BRANCH}
RUN cd gateway-mt && git reset --hard ${COMMIT}
WORKDIR gateway-mt

FROM base AS branch
ARG BRANCH
RUN git clone https://github.com/storj/gateway-mt.git --depth=1 --branch ${BRANCH}
WORKDIR gateway-mt

FROM ${SOURCE} AS github

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/gateway-mt.git
WORKDIR gateway-mt
RUN git fetch https://review.dev.storj.tools/storj/gateway-mt ${REF} && git checkout FETCH_HEAD

FROM base AS local
ARG PATH
WORKDIR /var/lib/storj/gateway-mt
COPY --chown=storj ${PATH} .

FROM --platform=$TARGETPLATFORM ${TYPE} AS binaries
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    go install -race ./cmd/...

FROM img.dev.storj.io/storjup/base:20250408-1 AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin