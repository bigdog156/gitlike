# Git Integration Guide

## Overview
The Todo CLI now provides seamless integration with your local Git repository, allowing you to:
- Automatically sync todo branches with Git branches
- Create Git commits alongside todo commits
- Track development progress with your code changes
- Maintain consistent workflow between todos and Git operations

## Setup

### 1. Initialize Git Integration
```bash
# Navigate to your Git repository
cd /path/to/your/git/repo

# Initialize todo Git integration
todo-cli git init
```

This will:
- Enable Git integration
- Detect your Git repository
- Sync existing Git branches with todo branches
- Set up auto-sync (optional)

### 2. Check Git Integration Status
```bash
todo-cli git status
```

Shows:
- Integration status (enabled/disabled)
- Current Git branch
- Repository path and remote URL
- Changed files count
- Last sync timestamp

## Core Features

### Automatic Branch Synchronization

#### Create Branch (Git + Todo)
```bash
# Creates both Git branch and todo branch
todo-cli git branch create feature-auth

# Equivalent to:
# git checkout -b feature-auth
# todo-cli branch create feature-auth
```

#### Switch Branch (Git + Todo)
```bash
# Switches both Git and todo branches
todo-cli git branch switch main

# Or use checkout
todo-cli git branch checkout feature-auth
```

#### List Branches (Synchronized View)
```bash
# Shows Git branches with todo counts
todo-cli git branch list
```

### Enhanced Todo Workflow

#### Regular Todo Operations
```bash
# Work with todos normally
todo-cli todo add "Implement authentication" -p high
todo-cli todo add "Add unit tests" -p medium

# Update todo status
todo-cli todo update 1 in-progress
todo-cli todo update 1 completed
```

#### Git-Integrated Commits
```bash
# Creates both todo commit AND Git commit
todo-cli git commit "Implement user authentication"

# With auto-push to Git remote
todo-cli git commit "Complete feature" --push
```

This creates:
1. Todo commit (tracks completed todos)
2. Git commit (commits your code changes)
3. Optional Git push to remote

### Auto-Sync Features

When Git integration is enabled with auto-sync:

#### Branch Switching
```bash
# Automatically syncs Git when switching todo branches
todo-cli branch switch feature-branch

# Behind the scenes:
# 1. Switches todo branch
# 2. Checks out Git branch
# 3. Syncs todo data
```

#### Manual Sync
```bash
# Manually sync todo branches with Git
todo-cli git sync
```

## Advanced Workflows

### Development Workflow Example
```bash
# 1. Start new feature
todo-cli git branch create feature-user-login

# 2. Plan your work
todo-cli todo add "Create login form" -p high
todo-cli todo add "Add validation" -p medium  
todo-cli todo add "Write tests" -p low

# 3. Work on tasks
todo-cli todo update 1 in-progress
# ... write code ...
todo-cli todo update 1 completed

# 4. Commit progress (todo + git)
todo-cli git commit "Implement login form UI"

# 5. Continue with next task
todo-cli todo update 2 in-progress
# ... write code ...
todo-cli todo update 2 completed

# 6. Final commit and push
todo-cli git commit "Add form validation" --push

# 7. Switch back to main
todo-cli git branch switch main
```

### Team Collaboration
```bash
# Pull latest changes from Git
git pull origin main

# Sync todo branches with Git
todo-cli git sync

# Work on assigned todos
todo-cli todo add "Fix reported bug" -p high
todo-cli todo update 1 completed

# Create integrated commit
todo-cli git commit "Fix authentication bug" --push
```

## Configuration

### Git Integration Settings
The integration settings are stored in your todo repository:

```json
{
  "git_integration": {
    "enabled": true,
    "auto_sync": true,
    "repo_path": "/path/to/git/repo",
    "remote_url": "https://github.com/user/repo.git",
    "auto_commit": false,
    "commit_template": "todo: {{.Message}}"
  }
}
```

### Enable/Disable Features
```bash
# Disable auto-sync (manual control)
# Edit ~/.tododata/repository.json and set "auto_sync": false

# Or reinitialize
todo-cli git init
```

## Command Reference

### Git Integration Commands
- `todo-cli git init` - Initialize Git integration
- `todo-cli git status` - Show integration status
- `todo-cli git sync` - Manual sync with Git
- `todo-cli git commit [msg]` - Create todo + Git commit
- `todo-cli git branch [action] [name]` - Git branch operations

### Enhanced Regular Commands (with Git integration)
- `todo-cli branch switch [name]` - Auto Git checkout if enabled
- `todo-cli commit create [msg]` - Regular todo commit only

## Benefits

### For Individual Developers
- **Unified workflow**: One tool for todos and Git
- **Progress tracking**: Correlate todos with actual commits
- **Branch isolation**: Todos automatically follow Git branches
- **Context switching**: Seamless branch switching

### for Teams
- **Synchronized development**: Todos travel with branches
- **Progress visibility**: Link todo completion to Git commits
- **Consistent workflow**: Standardized todo + Git operations
- **History tracking**: Complete audit trail of work

## Troubleshooting

### Git Repository Not Found
```bash
# Ensure you're in a Git repository
git status

# Initialize Git if needed
git init

# Then initialize todo Git integration
todo-cli git init
```

### Branch Sync Issues
```bash
# Manual resync
todo-cli git sync

# Check current status
todo-cli git status

# Verify Git branch
git branch --show-current
```

### Integration Disabled
```bash
# Re-enable integration
todo-cli git init

# Check status
todo-cli git status
```

## Tips

1. **Start with Git integration** from the beginning of projects
2. **Use descriptive commit messages** that include todo context
3. **Sync regularly** when working with teams
4. **Branch consistently** - let todos follow your Git branches
5. **Review integration status** periodically with `todo-cli git status`

The Git integration provides a seamless bridge between task management and version control, making your development workflow more organized and trackable!
