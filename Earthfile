VERSION 0.6
FROM golang:1.18
WORKDIR /go/storj-up

lint:
    FROM storjlabs/ci
    COPY . /go/storj-up
    WORKDIR /go/storj-up
    RUN check-copyright
    RUN check-large-files
    RUN check-imports -race ./...
    RUN check-atomic-align ./...
    RUN check-errs ./...
    RUN check-monkit ./...
    RUN staticcheck ./...
    RUN golangci-lint --build-tags mage -j=2 run
    RUN check-mod-tidy -mod .build/go.mod.orig

build:
    RUN ls -lah
    RUN pwd
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg/mod \
        go build -o build/ ./...
    SAVE ARTIFACT build/storj-up AS LOCAL build/storj-up

test:
   COPY . .
   RUN go install github.com/mfridman/tparse@36f80740879e24ba6695649290a240c5908ffcbb
   RUN mkdir build
   RUN --mount=type=cache,target=/root/.cache/go-build \
       --mount=type=cache,target=/go/pkg/mod \
       go test -json ./... | tee build/tests.json
   SAVE ARTIFACT build/tests.json AS LOCAL build/tests.json

integration:
   FROM earthly/dind:ubuntu
   RUN apt-get update && apt-get install -y docker-compose-plugin gcc
   RUN bash -c "curl --fail -L https://go.dev/dl/go1.17.12.linux-amd64.tar.gz | tar -C /usr/local -xz && cp /usr/local/go/bin/go /usr/local/bin/go"
   RUN go install github.com/rclone/rclone@v1.59.1
   RUN go install storj.io/storj/cmd/uplink@latest
   RUN go install storj.io/storjscan/cmd/storjscan@latest
   RUN go install github.com/rclone/rclone@v1.59.1
   COPY +build/storj-up /root/go/bin/storj-up
   ENV PATH=$PATH:/root/go/bin
   WORKDIR /test
   COPY ./test .
   WITH DOCKER
      RUN ./test.sh
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
