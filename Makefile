.PHONY: help build test coverage clean fmt deps

# Variables
GO := go
GOFLAGS := -v
BINARY_NAME := tfdocs-extra
BIN_DIR := bin
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Default target
help:
	@echo "Available targets:"
	@echo "  make build           - Build the binary"
	@echo "  make test            - Run tests"
	@echo "  make coverage        - Run tests with coverage report"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make fmt             - Format code"
	@echo "  make deps            - Download dependencies"
	@echo "  make all             - Build and test"

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd

test:
	$(GO) test -v ./...

coverage:
	$(GO) test -coverprofile=$(COVERAGE_FILE) ./...
	@echo "Coverage report generated: $(COVERAGE_FILE)"
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "HTML coverage report generated: $(COVERAGE_HTML)"

clean:
	@rm -rf $(BIN_DIR)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@$(GO) clean

fmt:
	$(GO) fmt ./...

deps:
	$(GO) mod download

all: fmt test build
	@echo "Build and test completed successfully!"
