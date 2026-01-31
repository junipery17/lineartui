BINARY_NAME=lineartui

BUILD_DIR=.

LDFLAGS=-ldflags "-s -w"

.PHONY: all
all: clean build

.PHONY: build
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

.PHONY: clean
clean:
	go clean
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: test
test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...
