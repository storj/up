FROM archlinux
RUN pacman -Sy --noconfirm which

FROM archlinux
COPY --from=0 /usr/bin/which /usr/bin/which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj
ADD pkg/recipe/entrypoint.sh /var/lib/storj/entrypoint.sh
ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
