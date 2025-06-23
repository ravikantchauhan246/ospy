#!/bin/bash

# Build script for Ospy

set -e

VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

echo "Building Ospy ${VERSION}..."
echo "Build time: ${BUILD_TIME}"
echo "Git commit: ${GIT_COMMIT}"

# Create dist directory
mkdir -p dist

# Build for current platform
echo "Building for current platform..."
go build -ldflags "${LDFLAGS}" -o dist/ospy ./cmd/ospy

# Build for multiple platforms
echo "Cross-compiling..."

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/ospy-linux-amd64 ./cmd/ospy

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/ospy-linux-arm64 ./cmd/ospy

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/ospy-darwin-amd64 ./cmd/ospy

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/ospy-darwin-arm64 ./cmd/ospy

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/ospy-windows-amd64.exe ./cmd/ospy

echo "Build complete! Binaries are in dist/"
ls -la dist/
