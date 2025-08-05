# Todo CLI Installation Guide

## Prerequisites
- Go 1.19 or higher installed on your system
- Git (optional, for cloning)

## Method 1: Build from Source

### Step 1: Get the source code
```bash
# Option A: If you have the code in a git repository
git clone <your-repo-url>
cd todocli

# Option B: Copy the source files manually to a new directory
mkdir todocli
cd todocli
# Copy all .go files and go.mod, go.sum
```

### Step 2: Install dependencies
```bash
go mod tidy
```

### Step 3: Build the application
```bash
# Build for current platform
go build -o todo-cli

# Build for specific platforms
# For Linux
GOOS=linux GOARCH=amd64 go build -o todo-cli-linux

# For Windows
GOOS=windows GOARCH=amd64 go build -o todo-cli.exe

# For macOS
GOOS=darwin GOARCH=amd64 go build -o todo-cli-macos
```

### Step 4: Install globally (optional)
```bash
# Copy to a directory in your PATH
sudo mv todo-cli /usr/local/bin/

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
GOOS=linux GOARCH=amd64 go build -o releases/todo-cli-linux-amd64
GOOS=linux GOARCH=arm64 go build -o releases/todo-cli-linux-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o releases/todo-cli-windows-amd64.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o releases/todo-cli-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o releases/todo-cli-darwin-arm64
```

### Step 2: Transfer and install
```bash
# Download/copy the appropriate binary for your system
# Make it executable (Linux/macOS)
chmod +x todo-cli-linux-amd64

# Move to PATH
sudo mv todo-cli-linux-amd64 /usr/local/bin/todo-cli
```

## Method 3: Go Install (if published)

If you publish your module to a Git repository:
```bash
go install github.com/yourusername/todocli@latest
```

## Usage After Installation

### Basic Commands
```bash
# Check installation
todo-cli --help

# Create a branch
todo-cli branch create feature-new

# Add a todo
todo-cli todo add "My first task" -d "Task description" -p high

# List todos
todo-cli todo list

# Update todo status
todo-cli todo update 1 completed

# Create commit
todo-cli commit create "Complete first task"

# Switch branches
todo-cli branch switch main

# Merge branches
todo-cli merge feature-new
```

### Data Location
- All data is stored in `~/.tododata/repository.json`
- Each user/device has independent data
- To migrate data, copy the `~/.tododata/` directory

## Troubleshooting

### Permission Issues
```bash
# If you get permission denied
chmod +x todo-cli
```

### Command not found
```bash
# Check if binary is in PATH
which todo-cli

# Add to PATH if needed
export PATH=$PATH:/path/to/todo-cli/directory
```

### Dependencies
```bash
# If build fails, ensure Go modules are working
go mod download
go mod verify
```

## Platform-Specific Notes

### Windows
- Use `todo-cli.exe` as the binary name
- Add to PATH via System Properties > Environment Variables

### Linux/Ubuntu
```bash
# Install via package manager (if you create .deb package)
sudo dpkg -i todo-cli.deb

# Or manual installation
sudo cp todo-cli /usr/local/bin/
```

### macOS
```bash
# For Apple Silicon Macs, use arm64 binary
# For Intel Macs, use amd64 binary

# Install via Homebrew (if you create a formula)
brew install todo-cli
```
