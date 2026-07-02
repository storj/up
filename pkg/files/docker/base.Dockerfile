FROM golang:1.26.4 AS storjup
COPY . /go/storj-up
WORKDIR /go/storj-up
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod  \
    go install

FROM ubuntu:noble AS final
RUN apt-get update
RUN apt-get -y install iproute2 ca-certificates curl jq unzip
RUN userdel -r ubuntu 2>/dev/null || true
RUN useradd storj --uid 1000 -d /storj && \
    mkdir -p /storj/shared && \
    mkdir -p /storj/bin && \
    chown storj. /storj
COPY --chown=storj identities /storj/identities
COPY --chown=storj --from=storjup /go/bin/storj-up /storj/bin/storj-up
USER storj
WORKDIR /storj
COPY pkg/recipe/entrypoint.sh /storj/entrypoint.sh
ENTRYPOINT ["/storj/entrypoint.sh"]
ENV PATH=$PATH:/storj/bin
LABEL org.opencontainers.image.source="https://github.com/storj/up"
