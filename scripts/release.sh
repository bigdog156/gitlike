#!/bin/bash

# GitLike CLI - Homebrew Release Script
# This script helps create releases for Homebrew publishing

set -e

VERSION=${1:-"1.0.0"}
REPO_URL="https://github.com/bigdog156/gitlike"

echo "🚀 Creating GitLike CLI release v${VERSION}"

# Clean previous builds
rm -rf releases
mkdir -p releases

echo "📦 Building binaries..."

# Build for macOS Intel
echo "  Building for macOS Intel (x86_64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o releases/gitlike-darwin-amd64

# Build for macOS Apple Silicon
echo "  Building for macOS Apple Silicon (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o releases/gitlike-darwin-arm64

# Build for Linux
echo "  Building for Linux (x86_64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o releases/gitlike-linux-amd64

# Build for Windows
echo "  Building for Windows (x86_64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o releases/gitlike-windows-amd64.exe

echo "✅ All binaries built successfully!"

# Generate SHA256 hashes
echo "🔐 Generating SHA256 hashes..."
cd releases
sha256sum gitlike-* > SHA256SUMS
echo ""
echo "SHA256 Hashes:"
cat SHA256SUMS
cd ..

# Create Homebrew formula template
echo "🍺 Creating Homebrew formula template..."

# Get SHA256 for macOS binaries
DARWIN_AMD64_SHA256=$(shasum -a 256 releases/gitlike-darwin-amd64 | cut -d' ' -f1)
DARWIN_ARM64_SHA256=$(shasum -a 256 releases/gitlike-darwin-arm64 | cut -d' ' -f1)

cat > homebrew-formula-template.rb << EOF
class Gitlike < Formula
  desc "GitLike CLI with Git-like workflow for developers"
  homepage "${REPO_URL}"
  version "${VERSION}"
  
  if Hardware::CPU.arm?
    url "${REPO_URL}/releases/download/v${VERSION}/gitlike-darwin-arm64"
    sha256 "${DARWIN_ARM64_SHA256}"
  else
    url "${REPO_URL}/releases/download/v${VERSION}/gitlike-darwin-amd64"
    sha256 "${DARWIN_AMD64_SHA256}"
  end

  def install
    bin.install "gitlike-darwin-arm64" => "gitlike" if Hardware::CPU.arm?
    bin.install "gitlike-darwin-amd64" => "gitlike" if Hardware::CPU.intel?
  end

  test do
    system "#{bin}/gitlike", "--help"
    assert_match "${VERSION}", shell_output("#{bin}/gitlike --version 2>&1")
  end
end
EOF

echo "✅ Homebrew formula template created: homebrew-formula-template.rb"

# Create GoReleaser config
echo "🎯 Creating GoReleaser configuration..."
cat > .goreleaser.yaml << EOF
project_name: gitlike

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -X main.version={{.Version}}

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
    format_overrides:
      - goos: windows
        format: zip

brews:
  - name: gitlike
    tap:
      owner: bigdog156
      name: homebrew-gitlike
    homepage: "${REPO_URL}"
    description: "GitLike CLI with Git-like workflow for developers"
    license: "MIT"
    test: |
      system "#{bin}/gitlike --help"
      assert_match version.to_s, shell_output("#{bin}/gitlike --version 2>&1")
    install: |
      bin.install "gitlike"

release:
  github:
    owner: bigdog156
    name: gitlike
EOF

echo "✅ GoReleaser configuration created: .goreleaser.yaml"

echo ""
echo "🎉 Release preparation complete!"
echo ""
echo "Next steps:"
echo "1. Create a Git tag: git tag v${VERSION} && git push origin v${VERSION}"
echo "2. Create a GitHub release and upload the binaries from 'releases/' directory"
echo "3. Create a repository named 'homebrew-gitlike' on GitHub"
echo "4. Copy the content from 'homebrew-formula-template.rb' to 'gitlike.rb' in your tap"
echo "5. Test with: brew install --build-from-source ./gitlike.rb"
echo ""
echo "Or use GoReleaser for automated releases:"
echo "  brew install goreleaser"
echo "  goreleaser release --rm-dist"
echo ""
echo "📁 Files created:"
echo "  - releases/ (directory with binaries)"
echo "  - homebrew-formula-template.rb"
echo "  - .goreleaser.yaml"
echo "  - releases/SHA256SUMS"
