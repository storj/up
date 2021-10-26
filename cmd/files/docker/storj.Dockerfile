ARG TYPE

FROM devbase AS base
ARG BRANCH=main
ARG REF=refs/changes/61/5961/9
ARG REPO=https://github.com/storj/storj

FROM base AS github
RUN git clone ${REPO} --depth=1 --branch ${BRANCH}
WORKDIR storj

FROM base AS gerrit
RUN git clone ${REPO} --branch ${BRANCH}
WORKDIR storj
RUN git fetch https://review.dev.storj.io/storj/storj ${REF} && git checkout FETCH_HEAD

FROM ${TYPE} AS binaries
RUN cd web/satellite && npm install && npm run build
RUN cd web/multinode && npm install && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build
RUN cd web/storagenode && npm install && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build
RUN env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh
RUN go install ./cmd/...

FROM ubuntubase AS final
COPY --from=binaries /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=binaries /var/lib/storj/storj/web/satellite/static /var/lib/storj/storj/web/satellite/static
COPY --from=binaries /var/lib/storj/storj/web/satellite/dist /var/lib/storj/storj/web/satellite/dist
COPY --from=binaries /var/lib/storj/storj/release/head/wasm /var/lib/storj/storj/web/satellite/static/wasm
COPY --from=binaries --chown=storj /var/lib/storj/identities /var/lib/storj/identities
COPY --from=binaries --chown=storj /var/lib/storj/entrypoint.sh /var/lib/storj/entrypoint.sh

ENTRYPOINT ["/var/lib/storj/entrypoint.sh"]
ENV PATH=$PATH:/var/lib/storj/go/bin
