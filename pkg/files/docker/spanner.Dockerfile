FROM img.dev.storj.io/storjup/base:20251110-1 AS binaries

FROM gcr.io/cloud-spanner-emulator/emulator:1.5.51 AS emulator

FROM debian:trixie-slim AS base

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

FROM base AS gcloud
ARG GCLOUD_VERSION=563.0.0
ARG TARGETARCH

RUN ARCH=$([ "$TARGETARCH" = "arm64" ] && echo "arm" || echo "x86_64") && \
    mkdir -p /opt && \
    curl -fsSL "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-${GCLOUD_VERSION}-linux-${ARCH}.tar.gz" \
    | tar -xz -C /opt && \
    mv /opt/google-cloud-sdk /opt/google-cloud-cli

FROM base AS final

COPY --from=gcloud /opt/google-cloud-cli /opt/google-cloud-cli
COPY --from=emulator /emulator_main /gateway_main /usr/local/bin/
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin

COPY pkg/recipe/startspanner.sh /var/lib/storj/startspanner.sh

ENV PATH=/opt/google-cloud-cli/bin:/var/lib/storj/go/bin:$PATH

ENTRYPOINT ["/var/lib/storj/startspanner.sh"]
