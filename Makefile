.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/ssowrap

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run
