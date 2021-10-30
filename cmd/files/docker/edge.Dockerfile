ARG TYPE

FROM ghcr.io/elek/storj-build:20211029-1 AS base

FROM base AS github
ARG BRANCH
RUN git clone --depth=1 https://github.com/storj/gateway-mt.git --branch ${BRANCH}
WORKDIR gateway-mt

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/gateway-mt.git
WORKDIR gateway-mt
RUN git fetch https://review.dev.storj.io/storj/gateway-mt ${REF} && git checkout FETCH_HEAD

FROM ${TYPE} AS binaries
RUN go install ./cmd/...

FROM ghcr.io/elek/storj-base:20211029-1 AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries /var/lib/storj/gateway-mt/pkg/linksharing/web /var/lib/storj/pkg/linksharing/web
COPY --from=binaries --chown=storj /var/lib/storj/identities /var/lib/storj/identities

