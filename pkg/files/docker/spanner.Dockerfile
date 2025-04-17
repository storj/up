FROM img.dev.storj.io/storjup/base:20250408-1 AS binaries

FROM gcr.io/cloud-spanner-emulator/emulator:1.5.31 AS emulator

FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:517.0.0-slim AS final

ADD pkg/recipe/startspanner.sh /var/lib/storj/startspanner.sh
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=emulator emulator_main .
COPY --from=emulator gateway_main .

ENTRYPOINT ["/var/lib/storj/startspanner.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
