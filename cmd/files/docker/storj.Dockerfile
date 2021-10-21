FROM ghcr.io/elek/storj-build
ARG BRANCH=v1.39.6
ARG REPO=https://github.com/storj/storj

RUN git clone ${REPO} --depth=1 --branch ${BRANCH}
RUN cd storj/web/satellite && npm install && npm run build
RUN cd storj/web/multinode && npm install && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build
RUN cd storj/web/storagenode && npm install && npm install && npm install @vue/cli-service && export PATH=$PATH:`pwd`/node_modules/.bin && npm run build
RUN cd storj && env env GO111MODULE=on GOOS=js GOARCH=wasm GOARM=6 -CGO_ENABLED=1 TAG=head scripts/build-wasm.sh
RUN cd storj && go install ./cmd/...

FROM ghcr.io/elek/storj-base

COPY --from=0 /var/lib/storj/go/bin /var/lib/storj/go/bin
COPY --from=0 /var/lib/storj/storj/web/satellite/static /var/lib/storj/storj/web/satellite/static
COPY --from=0 /var/lib/storj/storj/web/satellite/dist /var/lib/storj/storj/web/satellite/dist
COPY --from=0 /var/lib/storj/storj/web/satellite/static /var/lib/storj/storj/web/multinode/static
COPY --from=0 /var/lib/storj/storj/web/satellite/dist /var/lib/storj/storj/web/multinode/dist
COPY --from=0 /var/lib/storj/storj/web/satellite/static /var/lib/storj/storj/web/storagenode/static
COPY --from=0 /var/lib/storj/storj/web/satellite/dist /var/lib/storj/storj/web/storagenode/dist
COPY --from=0 /var/lib/storj/storj/release/head/wasm /var/lib/storj/storj/web/satellite/static/wasm
COPY --from=0 --chown=storj /var/lib/storj/identities /var/lib/storj/identities
