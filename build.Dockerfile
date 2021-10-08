FROM archlinux
ARG BRANCH=v1.39.6
ARG REPO=https://github.com/storj/storj
RUN pacman -Sy --noconfirm go git sudo npm make gcc which brotli
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj

USER storj
WORKDIR /var/lib/storj

RUN go install github.com/go-delve/delve/cmd/dlv@latest

#internal go mod chache
RUN git clone ${REPO} --depth=1 --branch ${BRANCH}  && \
    cd storj && \
    env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh && \
    go build ./cmd/... && \
    cd .. && \
    rm -rf storj

ADD pkg/recipe/entrypoint.sh /var/lib/storj/entrypoint.sh

ADD . /var/lib/storj/sjr
RUN cd /var/lib/storj/sjr/devrun && go install
COPY --chown=storj identities /var/lib/storj/identities
