# syntax=docker/dockerfile:1.3
FROM --platform=$TARGETPLATFORM ubuntu:noble as base
ARG TARGETPLATFORM
RUN apt-get update
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install curl && curl -sfL https://deb.nodesource.com/setup_22.x  | bash -
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install git sudo nodejs make gcc brotli g++
RUN echo ${TARGETPLATFORM} | sed 's/linux\///' | xargs -I PLATFORM curl --fail -L https://go.dev/dl/go1.26.4.linux-PLATFORM.tar.gz | tar -C /usr/local -xz && cp /usr/local/go/bin/go /usr/local/bin/go
ENV GOROOT=/usr/local/go
ENV GOBIN=/storj/bin
ENV PATH=$PATH:/storj/bin

RUN userdel -r ubuntu 2>/dev/null || true
RUN useradd storj --uid 1000 -d /storj && \
    mkdir -p /storj/shared && \
    mkdir -p /storj/.cache && \
    mkdir -p /storj/bin && \
    mkdir -p /storj/go/pkg/mod && \
    chown -R storj. /storj
USER storj
WORKDIR /storj
RUN --mount=type=cache,target=/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/storj/.cache/go-build,mode=777,uid=1000  \
    go install github.com/go-delve/delve/cmd/dlv@latest
COPY pkg/recipe/entrypoint.sh /storj/entrypoint.sh
LABEL org.opencontainers.image.source="https://github.com/storj/up"
