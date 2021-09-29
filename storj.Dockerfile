FROM archlinux
RUN pacman -Syu --noconfirm && pacman -S --noconfirm go git sudo npm make gcc which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj

USER storj
WORKDIR /var/lib/storj

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN git clone https://github.com/storj/storj --depth=1 --branch v1.39.6 && \
    cd storj/web/satellite && npm install && npm run build
RUN cd storj && env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh
RUN cd storj && go install ./cmd/...


ADD devrun /var/lib/storj/devrun
RUN cd /var/lib/storj/devrun && go install

FROM archlinux
RUN pacman -Syu --noconfirm which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj

COPY --from=0 /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=0 /var/lib/storj/storj/web/satellite/static /var/lib/storj/storj/web/satellite/static
COPY --from=0 /var/lib/storj/storj/web/satellite/dist /var/lib/storj/storj/web/satellite/dist
COPY --from=0 /var/lib/storj/storj/release/head/wasm /var/lib/storj/storj/web/satellite/static/wasm

ADD entrypoint.sh /var/lib/storj/entrypoint.sh
ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
