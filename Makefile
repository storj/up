.PHONY: lint-deps lint test-deps test build

lint-deps:
	go install github.com/storj/ci/...@045049f47d789130b7688358238c60ff11b31038
	go install honnef.co/go/tools/cmd/staticcheck@2025.1
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.1

lint:
	staticcheck ./...
	golangci-lint --build-tags mage -j=2 run
	check-copyright
	check-large-files
	check-imports -race ./...
	check-atomic-align ./...
	check-errs ./...
	check-monkit ./...
	check-mod-tidy

test-deps:
	go install github.com/mfridman/tparse@36f80740879e24ba6695649290a240c5908ffcbb
	sudo apt-get update && sudo apt-get install -y jq
	go install -race storj.io/storj/cmd/storagenode@v1.125.2
	go install -race storj.io/storj/cmd/satellite@v1.125.2
	go install -race storj.io/storj/cmd/versioncontrol@v1.125.2
	go install -race storj.io/edge/cmd/gateway-mt@v1.97.0
	go install -race storj.io/edge/cmd/linksharing@v1.97.0
	go install -race storj.io/edge/cmd/authservice@v1.97.0

test:
	./scripts/test.sh

build:
	go build -o ./build/storj-up .
