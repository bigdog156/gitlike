# Remote Repository Setup Guide

## Overview
The Todo CLI now supports remote repositories for synchronizing your todo data across multiple devices. This guide shows you how to set up and use remote functionality.

## Remote Types

### 1. HTTP Remote (Server-based)
Best for team collaboration and cloud synchronization.

### 2. File Remote (Local/Network file system)
Best for personal use across multiple devices with shared storage.

## Setup Instructions

### Method 1: HTTP Server Setup

#### Start the Todo CLI Server
```bash
# Build and run the server
cd /Users/lethachlam/Developer/todocli/server
go run main.go
# Server will start on http://localhost:8080
```

#### Configure Client
```bash
# Add HTTP remote
todo-cli remote add origin http://localhost:8080 -t http

# View configured remotes
todo-cli remote list
```

### Method 2: File-based Remote

#### Shared Directory Setup
```bash
# Create shared directory (could be on network drive, Dropbox, etc.)
mkdir -p ~/shared/todo-remote

# Add file remote
todo-cli remote add origin ~/shared/todo-remote/repository.json -t file
```

## Basic Remote Operations

### Push Changes
```bash
# Push to default remote (origin)
todo-cli push

# Push to specific remote
todo-cli push backup
```

### Pull Changes
```bash
# Pull and merge from default remote
todo-cli pull

# Pull from specific remote
todo-cli pull origin
```

### Fetch (Check Remote Changes)
```bash
# Fetch without merging
todo-cli fetch

# See what would be merged
todo-cli fetch origin
```

### Synchronize
```bash
# Pull then push (full sync)
todo-cli sync

# Sync with specific remote
todo-cli sync origin
```

## Advanced Features

### Remote Branch Switching
```bash
# Switch to branch and sync with remote
todo-cli branch switch feature-branch --sync

# Automatically pulls remote branch if it doesn't exist locally
```

### Multiple Remotes
```bash
# Add multiple remotes
todo-cli remote add origin http://main-server:8080 -t http
todo-cli remote add backup ~/backup/todo.json -t file
todo-cli remote add team http://team-server:8080 -t http

# Push to specific remote
todo-cli push backup
todo-cli push team
```

## Workflow Examples

### Individual Developer Workflow
```bash
# Set up remote
todo-cli remote add origin ~/Dropbox/todo-backup.json -t file

# Work normally
todo-cli todo add "New feature" -p high
todo-cli todo update 1 completed
todo-cli commit create "Complete new feature"

# Sync with remote
todo-cli push

# On another device
todo-cli pull  # Get latest changes
```

### Team Collaboration Workflow
```bash
# Team server setup
todo-cli remote add team http://team-server:8080 -t http

# Daily workflow
todo-cli pull team              # Get team updates
todo-cli branch create my-feature
todo-cli todo add "My task"
todo-cli todo update 1 completed
todo-cli commit create "My work"
todo-cli push team              # Share with team

# Switch to teammate's branch
todo-cli branch switch teammate-branch --sync
```

## Authentication

### HTTP Server Authentication
Set environment variables for authentication:
```bash
export TODO_CLI_USERNAME=your-username
export TODO_CLI_PASSWORD=your-password
```

### File System Permissions
Ensure read/write access to shared directories:
```bash
chmod 755 ~/shared/todo-remote
chmod 644 ~/shared/todo-remote/repository.json
```

## Troubleshooting

### Connection Issues
```bash
# Test server connection
curl http://localhost:8080/status

# Check remote configuration
todo-cli remote list
```

### Sync Conflicts
The system automatically merges changes, preferring newer updates:
- Todos with same ID: newer timestamp wins
- Branches: todos are merged
- Commits: duplicates are avoided

### Reset Remote
```bash
# Remove and re-add remote
todo-cli remote remove origin
todo-cli remote add origin http://new-server:8080 -t http
```

## Data Location
- Local data: `~/.tododata/repository.json`
- Server data: `server_repository.json` (in server directory)
- File remote: Specified path in remote URL

## Security Notes
- HTTP remotes support basic authentication
- File remotes rely on filesystem permissions
- Consider using HTTPS for production servers
- Backup your data regularly

## Server Deployment

### Local Development
```bash
cd server
go run main.go
```

### Production Deployment
```bash
# Build server binary
cd server
go build -o todo-server main.go

# Run with custom port
PORT=3000 ./todo-server

# Or deploy to cloud platform (Heroku, Railway, etc.)
```

The server provides a simple web interface at the root URL for monitoring.
