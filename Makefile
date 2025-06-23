# Makefile for Ospy

.PHONY: build clean test run deps cross-compile help

# Variables
VERSION ?= dev
BINARY_NAME = ospy
BUILD_DIR = dist
MAIN_PATH = ./cmd/ospy

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build with version information
build-versioned:
	@echo "Building $(BINARY_NAME) with version info..."
	@./scripts/build.sh

# Cross-compile for all platforms
cross-compile:
	@echo "Cross-compiling for all platforms..."
	@./scripts/build.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH) -config configs/config.yaml

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)

# Development build and run
dev: build run

# Show help
help:
	@echo "Available targets:"
	@echo "  build            - Build for current platform"
	@echo "  build-versioned  - Build with version information"
	@echo "  cross-compile    - Build for all platforms"
	@echo "  deps             - Install dependencies"
	@echo "  test             - Run tests"
	@echo "  run              - Run the application"
	@echo "  clean            - Clean build artifacts"
	@echo "  dev              - Build and run"
	@echo "  help             - Show this help"
