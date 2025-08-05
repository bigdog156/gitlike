# GitLike CLI Installation Guide

## Prerequisites
- Go 1.19 or higher installed on your system
- Git (optional, for cloning)

## Method 1: Build from Source

### Step 1: Get the source code
```bash
# Option A: If you have the code in a git repository
git clone <your-repo-url>
cd gitlike

# Option B: Copy the source files manually to a new directory
mkdir gitlike
cd gitlike
# Copy all .go files and go.mod, go.sum
```

### Step 2: Install dependencies
```bash
go mod tidy
```

### Step 3: Build the application
```bash
# Build for current platform
go build -o gitlike

# Build for specific platforms
# For Linux
GOOS=linux GOARCH=amd64 go build -o gitlike-linux

# For Windows
GOOS=windows GOARCH=amd64 go build -o gitlike.exe

# For macOS
GOOS=darwin GOARCH=amd64 go build -o gitlike-macos
```

### Step 4: Install globally (optional)
```bash
# Copy to a directory in your PATH
sudo mv gitlike /usr/local/bin/

# Or add current directory to PATH
echo 'export PATH=$PATH:$(pwd)' >> ~/.bashrc
source ~/.bashrc
```

## Method 2: Direct Binary Distribution

### Step 1: Create release binaries
```bash
# On your development machine, create binaries for different platforms
mkdir releases

# Linux
GOOS=linux GOARCH=amd64 go build -o releases/gitlike-linux-amd64
GOOS=linux GOARCH=arm64 go build -o releases/gitlike-linux-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o releases/gitlike-windows-amd64.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o releases/gitlike-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o releases/gitlike-darwin-arm64
```

### Step 2: Transfer and install
```bash
# Download/copy the appropriate binary for your system
# Make it executable (Linux/macOS)
chmod +x gitlike-linux-amd64

# Move to PATH
sudo mv gitlike-linux-amd64 /usr/local/bin/gitlike
```

## Method 3: Go Install (if published)

If you publish your module to a Git repository:
```bash
go install github.com/yourusername/gitlike@latest
```

## Usage After Installation

### Basic Commands
```bash
# Check installation
gitlike --help

# Create a branch
gitlike branch create feature-new

# Add a todo
gitlike todo add "My first task" -d "Task description" -p high

# List todos
gitlike todo list

# Update todo status
gitlike todo update 1 completed

# Create commit
gitlike commit create "Complete first task"

# Switch branches
gitlike branch switch main

# Merge branches
gitlike merge feature-new
```

### Data Location
- All data is stored in `~/.tododata/repository.json`
- Each user/device has independent data
- To migrate data, copy the `~/.tododata/` directory

## Troubleshooting

### Permission Issues
```bash
# If you get permission denied
chmod +x gitlike
```

### Command not found
```bash
# Check if binary is in PATH
which gitlike

# Add to PATH if needed
export PATH=$PATH:/path/to/gitlike/directory
```

### Dependencies
```bash
# If build fails, ensure Go modules are working
go mod download
go mod verify
```

## Platform-Specific Notes

### Windows
- Use `gitlike.exe` as the binary name
- Add to PATH via System Properties > Environment Variables

### Linux/Ubuntu
```bash
# Install via package manager (if you create .deb package)
sudo dpkg -i gitlike.deb

# Or manual installation
sudo cp gitlike /usr/local/bin/
```

### macOS
```bash
# For Apple Silicon Macs, use arm64 binary
# For Intel Macs, use amd64 binary

# Manual installation
sudo cp gitlike /usr/local/bin/

# Install via Homebrew (once published to a tap)
brew tap bigdog156/gitlike
brew install gitlike

# Or if published to Homebrew core
brew install gitlike
```
