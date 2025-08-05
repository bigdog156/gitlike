# Publishing GitLike CLI to Homebrew

## Overview
Homebrew is the most popular package manager for macOS and Linux. There are two main ways to publish your application:

1. **Official Homebrew Core** (very strict requirements)
2. **Personal Tap** (easier, recommended for most projects)

## Method 1: Personal Homebrew Tap (Recommended)

### Step 1: Create a GitHub Repository for Your Tap
```bash
# Create a new repository named "homebrew-<tapname>"
# For example: homebrew-gitlike
# Repository name MUST start with "homebrew-"
```

### Step 2: Create Release Binaries
First, create release binaries for different architectures:

```bash
# Create releases directory
mkdir -p releases

# Build for macOS Intel (x86_64)
GOOS=darwin GOARCH=amd64 go build -o releases/gitlike-darwin-amd64

# Build for macOS Apple Silicon (arm64)
GOOS=darwin GOARCH=arm64 go build -o releases/gitlike-darwin-arm64

# Build for Linux (optional)
GOOS=linux GOARCH=amd64 go build -o releases/gitlike-linux-amd64
```

### Step 3: Create GitHub Release
1. Tag your repository:
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

2. Create a GitHub release with your binaries
3. Upload the release binaries to the GitHub release

### Step 4: Create Homebrew Formula
Create a file called `gitlike.rb` in your homebrew tap repository:

```ruby
class Gitlike < Formula
  desc "GitLike CLI with Git-like workflow for developers"
  homepage "https://github.com/bigdog156/gitlike"
  version "1.0.0"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-darwin-arm64"
      sha256 "SHA256_HASH_FOR_ARM64_BINARY"
    else
      url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-darwin-amd64"
      sha256 "SHA256_HASH_FOR_AMD64_BINARY"
    end
  elsif OS.linux?
    url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/gitlike-linux-amd64"
    sha256 "SHA256_HASH_FOR_LINUX_BINARY"
  end

  def install
    bin.install "gitlike-darwin-arm64" => "gitlike" if Hardware::CPU.arm?
    bin.install "gitlike-darwin-amd64" => "gitlike" if Hardware::CPU.intel?
    bin.install "gitlike-linux-amd64" => "gitlike" if OS.linux?
  end

  test do
    system "#{bin}/gitlike", "--help"
  end
end
```

### Step 5: Generate SHA256 Hashes
```bash
# Generate SHA256 hashes for your binaries
shasum -a 256 releases/gitlike-darwin-amd64
shasum -a 256 releases/gitlike-darwin-arm64
shasum -a 256 releases/gitlike-linux-amd64
```

### Step 6: Publish Your Tap
1. Create repository: `homebrew-gitlike`
2. Add the formula file: `gitlike.rb`
3. Push to GitHub

### Step 7: Install Instructions for Users
```bash
# Add your tap
brew tap bigdog156/gitlike

# Install your application
brew install gitlike
```

## Method 2: Automated Release with GitHub Actions

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
        
    - name: Build binaries
      run: |
        # Build for different platforms
        GOOS=darwin GOARCH=amd64 go build -o todo-cli-darwin-amd64
        GOOS=darwin GOARCH=arm64 go build -o todo-cli-darwin-arm64
        GOOS=linux GOARCH=amd64 go build -o todo-cli-linux-amd64
        GOOS=windows GOARCH=amd64 go build -o todo-cli-windows-amd64.exe
        
    - name: Create Release
      uses: softprops/action-gh-release@v1
      with:
        files: |
          todo-cli-darwin-amd64
          todo-cli-darwin-arm64
          todo-cli-linux-amd64
          todo-cli-windows-amd64.exe
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Method 3: Using GoReleaser (Advanced)

Install GoReleaser:
```bash
brew install goreleaser
```

Create `.goreleaser.yaml`:
```yaml
project_name: todo-cli

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

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    format_overrides:
      - goos: windows
        format: zip

brews:
  - name: todo-cli
    tap:
      owner: bigdog156
      name: homebrew-todocli
    homepage: "https://github.com/bigdog156/gitlike"
    description: "Todo CLI with Git-like workflow for developers"
    license: "MIT"
    test: |
      system "#{bin}/todo-cli --help"
    install: |
      bin.install "todo-cli"

release:
  github:
    owner: bigdog156
    name: gitlike
```

Then release with:
```bash
goreleaser release --rm-dist
```

## Step-by-Step Quick Start

1. **Create the tap repository:**
   ```bash
   # On GitHub, create: homebrew-todocli
   ```

2. **Build and release your binaries:**
   ```bash
   # Tag your main repository
   git tag v1.0.0
   git push origin v1.0.0
   
   # Build binaries
   GOOS=darwin GOARCH=amd64 go build -o todo-cli-darwin-amd64
   GOOS=darwin GOARCH=arm64 go build -o todo-cli-darwin-arm64
   
   # Get SHA256 hashes
   shasum -a 256 todo-cli-darwin-*
   ```

3. **Create the formula:**
   ```bash
   # In your homebrew-todocli repository
   cat > todo-cli.rb << 'EOF'
   class TodoCli < Formula
     desc "Todo CLI with Git-like workflow for developers"
     homepage "https://github.com/bigdog156/gitlike"
     version "1.0.0"
     
     if Hardware::CPU.arm?
       url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/todo-cli-darwin-arm64"
       sha256 "YOUR_ARM64_SHA256_HERE"
     else
       url "https://github.com/bigdog156/gitlike/releases/download/v1.0.0/todo-cli-darwin-amd64"
       sha256 "YOUR_AMD64_SHA256_HERE"
     end

     def install
       bin.install "todo-cli-darwin-arm64" => "todo-cli" if Hardware::CPU.arm?
       bin.install "todo-cli-darwin-amd64" => "todo-cli" if Hardware::CPU.intel?
     end

     test do
       system "#{bin}/todo-cli", "--help"
     end
   end
   EOF
   ```

4. **Test your formula:**
   ```bash
   brew install --build-from-source ./todo-cli.rb
   ```

## User Installation Instructions

Once published, users can install with:
```bash
# Add your tap
brew tap bigdog156/todocli

# Install your CLI
brew install todo-cli

# Use it
todo-cli --help
```

## Tips for Success

1. **Choose meaningful names** - Repository should be `homebrew-<name>`
2. **Test thoroughly** - Test on both Intel and Apple Silicon Macs
3. **Keep it updated** - Update formula when you release new versions
4. **Add good documentation** - Clear description and homepage
5. **Follow conventions** - Use standard Homebrew formula patterns

## Alternative: Official Homebrew Core

For official Homebrew core (much harder):
1. Your project must be popular and well-maintained
2. Must follow strict guidelines
3. Submit PR to homebrew/homebrew-core
4. Very lengthy review process

## Next Steps

1. Create your `homebrew-todocli` repository
2. Create a release with binaries
3. Write the formula
4. Test it locally
5. Share installation instructions with users

Your todo CLI will then be installable with simple `brew install` commands!
