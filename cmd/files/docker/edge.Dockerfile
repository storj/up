ARG TYPE

FROM devbase AS base
ARG BRANCH=main
ARG REF=refs/changes/96/5696/31
ARG REPO=https://github.com/storj/gateway-mt

FROM base AS github
RUN git clone --depth=1 ${REPO} --branch ${BRANCH}
WORKDIR gateway-mt

FROM base AS gerrit
RUN git clone ${REPO} --branch ${BRANCH}
WORKDIR gateway-mt
RUN git fetch https://review.dev.storj.io/storj/gateway-mt ${REF} && git checkout FETCH_HEAD

FROM ${TYPE} AS binaries
RUN go install ./cmd/...

FROM ubuntubase AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries /var/lib/storj/gateway-mt/pkg/linksharing/web /var/lib/storj/pkg/linksharing/web
COPY --from=binaries --chown=storj /var/lib/storj/identities /var/lib/storj/identities

