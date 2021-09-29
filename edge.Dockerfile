FROM archlinux
RUN pacman -Syu --noconfirm && pacman -S --noconfirm go git sudo npm make gcc which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj

USER storj
WORKDIR /var/lib/storj

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN git clone --depth=1 https://github.com/storj/gateway-mt --branch v1.14.0 && \
   cd gateway-mt && go install ./cmd/...

ADD devrun /var/lib/storj/devrun
RUN cd /var/lib/storj/devrun && go install


FROM elek/storj

FROM archlinux
RUN pacman -Syu --noconfirm which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj

COPY --from=0 /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=0 /var/lib/storj/gateway-mt/pkg/linksharing/web /var/lib/storj/pkg/linksharing/web
COPY --from=1 /var/lib/storj/go/bin/identity /var/lib/storj/go/bin/identity
COPY --from=1 /var/lib/storj/go/bin/uplink /var/lib/storj/go/bin/uplink

ADD entrypoint.sh /var/lib/storj/entrypoint.sh
ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
