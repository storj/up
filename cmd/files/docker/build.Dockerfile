ARG TYPE

FROM ubuntu:21.04 as base
RUN apt-get update
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install curl && curl -sfL https://deb.nodesource.com/setup_16.x  | bash -
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install golang git sudo nodejs make gcc brotli
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj

USER storj
WORKDIR /var/lib/storj

RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM base AS github
ARG BRANCH
RUN git clone https://github.com/storj/storj.git --depth=1 --branch ${BRANCH}
WORKDIR storj

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/storj.git
WORKDIR storj
RUN git fetch https://review.dev.storj.io/storj/storj ${REF} && git checkout FETCH_HEAD

FROM ${TYPE} AS binaries
RUN env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh && \
    go build ./cmd/... && \
    cd .. && \
    rm -rf storj
WORKDIR ../

FROM binaries AS final
ADD pkg/recipe/entrypoint.sh /var/lib/storj/entrypoint.sh

ADD . /var/lib/storj/sjr
RUN cd /var/lib/storj/sjr/devrun && go install
COPY --chown=storj identities /var/lib/storj/identities
