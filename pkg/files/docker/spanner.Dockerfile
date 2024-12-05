FROM img.dev.storj.io/storjup/base:20240509-1 AS binaries

FROM google/cloud-sdk:502.0.0-slim AS final
RUN apt-get install -y google-cloud-cli-spanner-emulator

ADD pkg/recipe/startspanner.sh /var/lib/storj/startspanner.sh
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin

ENTRYPOINT ["/var/lib/storj/startspanner.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
