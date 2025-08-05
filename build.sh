#!/bin/bash

# Todo CLI Build Script
# Creates binaries for multiple platforms

set -e

VERSION=${1:-"v1.0.0"}
APP_NAME="todo-cli"
BUILD_DIR="releases"

echo "Building ${APP_NAME} ${VERSION}..."

# Clean and create build directory
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build for multiple platforms
echo "Building binaries..."

# Linux
echo "  - Linux (amd64)"
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ${BUILD_DIR}/${APP_NAME}-linux-amd64

echo "  - Linux (arm64)"
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o ${BUILD_DIR}/${APP_NAME}-linux-arm64

# Windows
echo "  - Windows (amd64)"
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe

# macOS
echo "  - macOS (amd64)"
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ${BUILD_DIR}/${APP_NAME}-darwin-amd64

echo "  - macOS (arm64)"
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o ${BUILD_DIR}/${APP_NAME}-darwin-arm64

# Create checksums
echo "Creating checksums..."
cd ${BUILD_DIR}
sha256sum * > checksums.txt
cd ..

echo "Build complete! Files in ${BUILD_DIR}:"
ls -la ${BUILD_DIR}/

echo ""
echo "To install locally:"
echo "  chmod +x ${BUILD_DIR}/${APP_NAME}-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/')"
echo "  sudo mv ${BUILD_DIR}/${APP_NAME}-* /usr/local/bin/${APP_NAME}"
