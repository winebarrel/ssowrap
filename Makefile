.PHONY: all
all: vet build

.PHONY: build
build:
	CGO_ENABLED=0 go build $(BUILD_OPTS) ./cmd/ssowrap

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run
