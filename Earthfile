VERSION 0.8
FROM golang:1.24.3
WORKDIR /go/storj-up

lint-deps:
    RUN go install github.com/storj/ci/...@045049f47d789130b7688358238c60ff11b31038
    RUN go install honnef.co/go/tools/cmd/staticcheck@2025.1
    RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.1

lint:
    FROM +lint-deps
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

build-app-deps:
    # Download deps before copying code.
    COPY go.mod go.sum ./pkg/config/gen/go.mod ./pkg/config/gen/go.sum .
    RUN go mod download
    # Output these back in case go mod download changes them.
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum
    SAVE ARTIFACT ./pkg/config/gen/go.mod AS LOCAL ./pkg/config/gen/go.mod
    SAVE ARTIFACT ./pkg/config/gen/go.sum AS LOCAL ./pkg/config/gen/go.sum

build-app:
    FROM +build-app-deps
    # Copy and build code.
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        go build -o build/ ./...
    SAVE ARTIFACT build/storj-up AS LOCAL build/storj-up

test:
   RUN go install github.com/mfridman/tparse@36f80740879e24ba6695649290a240c5908ffcbb
   RUN apt-get update && apt-get install -y jq
   RUN go install -race storj.io/storj/cmd/storagenode@v1.125.2
   RUN go install -race storj.io/storj/cmd/satellite@v1.125.2
   RUN go install -race storj.io/storj/cmd/versioncontrol@v1.125.2
   RUN go install -race storj.io/edge/cmd/gateway-mt@v1.97.0
   RUN go install -race storj.io/edge/cmd/linksharing@v1.97.0
   RUN go install -race storj.io/edge/cmd/authservice@v1.97.0
   RUN mkdir build
   COPY . .
   RUN --mount=type=cache,target=/root/.cache/go-build \
       --mount=type=cache,target=/go/pkg/mod \
       ./scripts/test.sh
   SAVE ARTIFACT build/tests.json AS LOCAL build/tests.json

integration-all:
   BUILD ./test/uplink+test
   BUILD ./test/spanner/uplink+test
   BUILD ./test/edge+test
   BUILD ./test/storjscan+test

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
