.PHONY: all
all: vet build

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags '-X main.version=$(VERSION)' ./cmd/ssowrap

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: image
image:
	docker build -t ssowrap .
