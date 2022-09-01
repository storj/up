FROM --platform=$TARGETPLATFORM ubuntu:22.04 AS final
RUN apt-get update
RUN apt-get -y install iproute2 ca-certificates
RUN useradd storj --uid 1000 -d /var/lib/storj && mkdir -p /var/lib/storj/shared && chown storj. /var/lib/storj
USER storj
WORKDIR /var/lib/storj
ADD pkg/recipe/entrypoint.sh /var/lib/storj/entrypoint.sh
ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
