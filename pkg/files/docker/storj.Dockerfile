# syntax=docker/dockerfile:1.3
ARG TYPE
ARG SOURCE
FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/build:20240911-1  AS base

ARG SKIP_FRONTEND_BUILD

FROM base AS commit
ARG BRANCH
ARG COMMIT
RUN git clone https://github.com/storj/storj.git --branch ${BRANCH}
RUN cd storj && git reset --hard ${COMMIT}
WORKDIR storj

FROM base AS branch
ARG BRANCH
RUN git clone https://github.com/storj/storj.git --depth=1 --branch ${BRANCH}
WORKDIR storj

FROM ${SOURCE} AS github

FROM base AS gerrit
ARG REF
RUN git clone https://github.com/storj/storj.git
WORKDIR storj
RUN git fetch https://review.dev.storj.io/storj/storj ${REF} && git checkout FETCH_HEAD

FROM base AS local
ARG PATH
WORKDIR /var/lib/storj/storj
COPY --chown=storj ${PATH} .

FROM ${TYPE} AS binaries
RUN if [ -z "$SKIP_FRONTEND_BUILD" ] ; then cd web/satellite && npm install && npm run build && npm run build-vuetify ; fi
RUN if [ -z "$SKIP_FRONTEND_BUILD" ] ; then cd web/multinode && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build ; fi
RUN if [ -z "$SKIP_FRONTEND_BUILD" ] ; then cd web/storagenode && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build ; fi
RUN if [ -z "$SKIP_FRONTEND_BUILD" ] ; then cd satellite/admin/ui && npm install && npm run build ; fi
RUN if [ -z "$SKIP_FRONTEND_BUILD" ] ; then cd satellite/admin/back-office/ui && npm install && npm run build ; fi
RUN if [ -z "$SKIP_FRONTEND_BUILD" ] ; then env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh ; fi

RUN --mount=type=cache,target=/var/lib/storj/go/pkg/mod,mode=777,uid=1000 \
    --mount=type=cache,target=/var/lib/storj/.cache/go-build,mode=777,uid=1000 \
    go install -race ./cmd/... \
    && go install -ldflags \
    "-X storj.io/common/version.buildRelease=false  \
    -X storj.io/common/version.buildVersion=v0.0.0  \
    -X storj.io/common/version.buildTimestamp=0"  \
    ./cmd/storagenode/...

FROM --platform=$TARGETPLATFORM img.dev.storj.io/storjup/base:20240509-1 AS final
ENV STORJ_ADMIN_STATIC_DIR=/var/lib/storj/storj/satellite/admin/ui/build
ENV STORJ_CONSOLE_STATIC_DIR=/var/lib/storj/storj/web/satellite/
ENV STORJ_MAIL_TEMPLATE_PATH=/var/lib/storj/storj/web/satellite/static/emails
ENV STORJ_STORAGENODE_CONSOLE_STATIC_DIR=/var/lib/storj/web/storagenode
# copy objects. the '[]' are to avoid the COPY command from failing if the files are missing.
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries /var/lib/storj/storj/web/satellite/stati[c] /var/lib/storj/storj/web/satellite/static
COPY --from=binaries /var/lib/storj/storj/web/satellite/dis[t] /var/lib/storj/storj/web/satellite/dist
COPY --from=binaries /var/lib/storj/storj/web/satellite/dist_vuetify_poc /var/lib/storj/storj/web/satellite/dist_vuetify_poc
COPY --from=binaries /var/lib/storj/storj/satellite/admin/ui/build /var/lib/storj/storj/satellite/admin/ui/build
COPY --from=binaries /var/lib/storj/storj/satellite/admin/back-office/ui/build /var/lib/storj/storj/satellite/admin/back-office/ui/build
COPY --from=binaries /var/lib/storj/storj/web/storagenode/stati[c] /var/lib/storj/web/storagenode/static
COPY --from=binaries /var/lib/storj/storj/web/storagenode/dis[t] /var/lib/storj/web/storagenode/dist
COPY --from=binaries /var/lib/storj/storj/web/multinode/stati[c] /var/lib/storj/web/multinode/static
COPY --from=binaries /var/lib/storj/storj/web/multinode/dis[t] /var/lib/storj/web/multinode/dist
COPY --from=binaries /var/lib/storj/storj/releas[e]/head/wasm /var/lib/storj/storj/web/satellite/static/wasm
COPY --from=binaries --chown=storj /var/lib/storj/entrypoint.sh /var/lib/storj/entrypoint.sh

ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
