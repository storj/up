FROM ghcr.io/elek/storj-build

ARG BRANCH=v1.14.0
ARG REPO=https://github.com/storj/gateway-mt
RUN git clone --depth=1 ${REPO} --branch ${BRANCH} && \
   cd gateway-mt && go install ./cmd/...

FROM ghcr.io/elek/storj-base

COPY --from=0 /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=0 /var/lib/storj/gateway-mt/pkg/linksharing/web /var/lib/storj/pkg/linksharing/web
COPY --from=0 --chown=storj /var/lib/storj/identities /var/lib/storj/identities

