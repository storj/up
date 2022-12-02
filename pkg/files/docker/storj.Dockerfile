# syntax=docker/dockerfile:1.3
ARG TYPE
FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/build:20220901-2  AS base

FROM base AS github
ARG BRANCH
RUN git clone https://github.com/storj/storj.git --depth=1 --branch ${BRANCH}
WORKDIR storj

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/storj.git
WORKDIR storj
RUN git fetch https://review.dev.storj.io/storj/storj ${REF} && git checkout FETCH_HEAD

FROM ${TYPE} AS binaries
RUN cd web/satellite && npm install && npm run build
RUN cd web/multinode && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build
RUN cd web/storagenode && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build
RUN cd satellite/admin/ui && npm install && npm run build
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh
RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    go install ./cmd/...

FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/base:20220901-3 AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries /var/lib/storj/storj/web/satellite/static /var/lib/storj/storj/web/satellite/static
COPY --from=binaries /var/lib/storj/storj/web/satellite/dist /var/lib/storj/storj/web/satellite/dist
COPY --from=binaries /var/lib/storj/storj/satellite/admin/ui/build /var/lib/storj/storj/satellite/admin/ui/build
COPY --from=binaries /var/lib/storj/storj/web/storagenode/static /var/lib/storj/web/storagenode/static
COPY --from=binaries /var/lib/storj/storj/web/storagenode/dist /var/lib/storj/web/storagenode/dist
COPY --from=binaries /var/lib/storj/storj/web/multinode/static /var/lib/storj/web/multinode/static
COPY --from=binaries /var/lib/storj/storj/web/multinode/dist /var/lib/storj/web/multinode/dist
COPY --from=binaries /var/lib/storj/storj/release/head/wasm /var/lib/storj/storj/web/satellite/static/wasm
COPY --from=binaries --chown=storj /var/lib/storj/entrypoint.sh /var/lib/storj/entrypoint.sh

ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
