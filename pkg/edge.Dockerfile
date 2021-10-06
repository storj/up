FROM ghcr.io/elek/storj-build

ARG BRANCH=v1.14.0
ARG REPO=https://github.com/storj/gateway-mt
RUN git clone --depth=1 ${REPO} --branch ${BRANCH} && \
   cd gateway-mt && go install ./cmd/...

FROM archlinux
RUN pacman -Syu --noconfirm which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj

COPY --from=0 /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=0 /var/lib/storj/gateway-mt/pkg/linksharing/web /var/lib/storj/pkg/linksharing/web
COPY --from=0 --chown=storj /var/lib/storj/identities /var/lib/storj/identities
COPY --from=0 --chown=storj /var/lib/storj/entrypoint.sh /var/lib/storj/entrypoint.sh

ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
