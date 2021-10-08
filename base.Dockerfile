FROM archlinux
RUN pacman -Syu --noconfirm which
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj
ADD pkg/entrypoint.sh /var/lib/storj/entrypoint
ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
