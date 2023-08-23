VERSION 0.6
FROM golang:1.19
WORKDIR /go/storj-up

lint:
    WORKDIR /go/storj-up
    RUN go install github.com/storj/ci/...@5286f52
    RUN go install honnef.co/go/tools/cmd/staticcheck@2022.1.3
    RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.3
    COPY . /go/storj-up
    RUN staticcheck ./...
    RUN golangci-lint --build-tags mage -j=2 run
    RUN check-copyright
    RUN check-large-files
    RUN check-imports -race ./...
    RUN check-atomic-align ./...
    RUN check-errs ./...
    RUN check-monkit ./...
    RUN check-mod-tidy

build:
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        go build -o build/ ./...
    SAVE ARTIFACT build/storj-up AS LOCAL build/storj-up

test:
   RUN go install github.com/mfridman/tparse@36f80740879e24ba6695649290a240c5908ffcbb
   RUN apt-get update && apt-get install -y jq
   RUN go install storj.io/storj/cmd/storagenode@v1.65.1
   RUN go install storj.io/storj/cmd/satellite@v1.65.1
   RUN go install storj.io/gateway-mt/cmd/gateway-mt@v1.39.0
   RUN go install storj.io/gateway-mt/cmd/linksharing@v1.39.0
   RUN go install storj.io/gateway-mt/cmd/authservice@v1.39.0
   RUN mkdir build
   COPY . .
   RUN --mount=type=cache,target=/root/.cache/go-build \
       --mount=type=cache,target=/go/pkg/mod \
       ./scripts/test.sh
   SAVE ARTIFACT build/tests.json AS LOCAL build/tests.json

integration:
   FROM earthly/dind:ubuntu
   RUN apt-get update && apt-get install -y golang-go docker-compose-plugin gcc
   ENV PATH=$PATH:/root/go/bin
   COPY . .
   WORKDIR /test
   WITH DOCKER --pull img.dev.storj.io/storjup/storj:1.83.2 \
               --pull img.dev.storj.io/storjup/edge:1.58.1 \
               --pull cockroachdb/cockroach \
               --pull redis:6.0.9 \
               --pull img.dev.storj.io/storjup/storjscan:1.0.3 \
               --pull ethereum/client-go
      RUN ./test.sh && \
          ./test-edge.sh && \
          ./test-storjscan.sh
   END

check-format:
   COPY . .
   RUN mkdir build
   RUN bash -c '[[ $(git status --short) == "" ]] || (echo "Before formatting, please commit all your work!!! (Formatter will format only last commit)" && exit -1)'
   RUN git show --name-only --pretty=format: | grep ".go" | xargs -n1 gofmt -s -w
   RUN git diff > build/format.patch
   SAVE ARTIFACT build/format.patch

check-format-all:
   RUN go install mvdan.cc/gofumpt@v0.3.1
   COPY . /go/storj-up
   WORKDIR /go/storj-up
   RUN bash -c 'find -name "*.go" | xargs -n1 gofmt -s -w'
   RUN bash -c 'find -name "*.go" | xargs -n1 gofumpt -s -w'
   RUN mkdir -p build
   RUN git diff > build/format.patch
   SAVE ARTIFACT build/format.patch

format:
   LOCALLY
   COPY +check-format/format.patch build/format.patch
   RUN git apply --allow-empty build/format.patch
   RUN git status

format-all:
   LOCALLY
   COPY +check-format-all/format.patch build/format.patch
   RUN git apply --allow-empty build/format.patch
   RUN git status
