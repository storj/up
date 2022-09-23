# syntax=docker/dockerfile:1.3
ARG TYPE
FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/build:20220803-2 AS base

FROM base AS github
ARG BRANCH
RUN git clone --depth=1 https://github.com/storj/storjscan.git --branch ${BRANCH}
WORKDIR storjscan

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/storjscan.git
WORKDIR storjscan
RUN git fetch https://review.dev.storj.io/storj/storjscan ${REF} && git checkout FETCH_HEAD

FROM --platform=$TARGETPLATFORM ${TYPE} AS binaries
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    go install ./cmd/...

FROM img.dev.storj.io/storjup/base:20220901-3 AS final

COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries --chown=storj /var/lib/storj/entrypoint.sh /var/lib/storj/entrypoint.sh

ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin