# syntax=docker/dockerfile:1.3
ARG TYPE
FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/build:20221220-2 AS base

FROM base AS github
ARG BRANCH
RUN git clone --depth=1 https://github.com/storj/gateway-mt.git --branch ${BRANCH}
WORKDIR gateway-mt

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/gateway-mt.git
WORKDIR gateway-mt
RUN git fetch https://review.dev.storj.io/storj/gateway-mt ${REF} && git checkout FETCH_HEAD

FROM --platform=$TARGETPLATFORM ${TYPE} AS binaries
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    go install ./cmd/...

FROM img.dev.storj.io/storjup/base:20221220-2 AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries /var/lib/storj/gateway-mt/pkg/linksharing/web /var/lib/storj/pkg/linksharing/web


