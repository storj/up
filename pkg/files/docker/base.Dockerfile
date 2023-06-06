FROM --platform=$TARGETPLATFORM golang:1.19 AS storjup
COPY . /go/storj-up
WORKDIR /go/storj-up
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod  \
    go install

FROM --platform=$TARGETPLATFORM ubuntu:22.04 AS final
RUN apt-get update
RUN apt-get -y install iproute2 ca-certificates
RUN useradd storj --uid 1000 -d /var/lib/storj && \
    mkdir -p /var/lib/storj/shared && \
    mkdir -p /var/lib/storj/go/bin && \
    chown storj. /var/lib/storj
COPY --chown=storj identities /var/lib/storj/identities
COPY --chown=storj --from=storjup /go/bin/storj-up /var/lib/storj/go/bin/storj-up
USER storj
WORKDIR /var/lib/storj
ADD pkg/recipe/entrypoint.sh /var/lib/storj/entrypoint.sh
ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
