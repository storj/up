FROM gcr.io/cloud-spanner-emulator/emulator:1.5.52 AS emulator

FROM debian:trixie-slim AS final

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/*

COPY --from=emulator /emulator_main /gateway_main /usr/local/bin/

COPY pkg/recipe/startspanner.sh /var/lib/storj/startspanner.sh

ENTRYPOINT ["/var/lib/storj/startspanner.sh"]
