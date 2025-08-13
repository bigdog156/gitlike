# GitLike CLI

> A Git-like CLI application with integrated todo tracking for task-driven development

GitLike mirrors Git's exact syntax while adding powerful todo management that automatically tracks which tasks are associated with commits. Perfect for developers who want to maintain clear task context throughout their Git workflow.

## 🚀 Features

- **🎯 Git-Identical Syntax**: All Git commands work exactly the same (`init`, `add`, `commit`, `status`, `log`, `checkout`, etc.)
- **📝 Automatic Todo Tracking**: Commits automatically link to active todos
- **🌿 Branch-Specific Todos**: Each branch maintains its own todo list
- **⚡ Active Task Workflow**: Start/stop working on todos with automatic commit linking
- **📊 Rich Commit History**: View commits with associated todo context
- **🔄 Git Integration**: Full compatibility with existing Git repositories

## 📦 Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd gitlike

# Build the application  
go build -o gitlike .

# Move to PATH (optional)
sudo mv gitlike /usr/local/bin/
```

## 🎯 Quick Start

```bash
# Initialize a GitLike repository (same as git init)
gitlike init

# Create a todo task
gitlike todo create "Implement user authentication" -d "Add JWT auth" -p high

# Start working on the todo (makes it active)
gitlike todo start 1

# Work normally with Git-like commands
echo "# Auth System" > auth.js
gitlike add auth.js
gitlike commit -m "Add authentication scaffolding"
# ✅ Git commit successful (working on #1: Implement user authentication)

# Mark todo as completed
gitlike todo done 1

# View commit history with todo tracking
gitlike log
```

## 📋 Git-Like Commands

### Core Git Commands (identical syntax)

| Command | Description | Example |
|---------|-------------|---------|
| `gitlike init` | Initialize repository | `gitlike init` |
| `gitlike add <files>` | Stage changes | `gitlike add . or gitlike add file.js` |
| `gitlike commit -m "msg"` | Commit with message | `gitlike commit -m "Add feature"` |
| `gitlike status` | Show working tree status | `gitlike status` |
| `gitlike log` | Show commit history | `gitlike log --oneline` |
| `gitlike checkout <branch>` | Switch branches | `gitlike checkout main` |
| `gitlike checkout -b <branch>` | Create & switch branch | `gitlike checkout -b feature-auth` |
| `gitlike push [remote] [branch]` | Push to remote | `gitlike push origin main` |
| `gitlike pull [remote] [branch]` | Pull from remote | `gitlike pull origin main` |

### Todo Commands (GitLike extensions)

| Command | Description | Example |
|---------|-------------|---------|
| `gitlike todo create "title"` | Create new todo | `gitlike todo create "Add login" -p high` |
| `gitlike todo done <id>` | Mark todo complete | `gitlike todo done 1` |
| `gitlike todo list` | List branch todos | `gitlike todo list` |
| `gitlike todo start <id>` | Start working on todo | `gitlike todo start 1` |
| `gitlike todo stop` | Stop current todo | `gitlike todo stop` |
| `gitlike todo active` | Show active todo | `gitlike todo active` |

### Advanced Commands

| Command | Description |
|---------|-------------|
| `gitlike branch` | Branch management |
| `gitlike git init` | Git integration setup |
| `gitlike git status` | Git integration status |

## 🔄 Task-Driven Workflow

GitLike introduces a task-driven development approach:

### 1. Create & Start Todo
```bash
gitlike todo create "Implement user authentication"
gitlike todo start 1
# 🚀 Started working on todo #1
# 💡 All commits will be automatically linked to this task!
```

### 2. Work Normally
```bash
# Work with regular Git-like commands
gitlike add auth.js
gitlike commit -m "Add JWT authentication"
# ✅ Git commit successful (working on #1: Implement user authentication)
```

### 3. Complete Todo
```bash
gitlike todo done 1
# ✅ Todo #1 marked as completed
# 💡 Ready to commit your work:
#    gitlike add .
#    gitlike commit -m "Complete todo #1"
```

### 4. View History
```bash
gitlike log
# commit abc1234
# Author: developer
# Date: Wed Aug 13 13:26:39 2025
#
#     Add JWT authentication
#
#     Todos included: #1
#     Active todo: #1
```

## 🌿 Branch-Specific Todo Management

Each branch maintains its own todo list:

```bash
# Create feature branch
gitlike checkout -b feature-login

# Branch-specific todos
gitlike todo create "Add login form"
gitlike todo create "Add form validation"

# Switch branches to see different todos
gitlike checkout main
gitlike todo list  # Shows main branch todos

gitlike checkout feature-login  
gitlike todo list  # Shows feature branch todos
```

## 📊 Todo Status & Flags

### Todo Flags
- `-d, --description`: Add detailed description
- `-p, --priority`: Set priority (`low`, `medium`, `high`)

### Todo Status Indicators
- `⏳` Pending
- `🔄` In Progress (active)
- `✅` Completed

### Example
```bash
gitlike todo create "Implement OAuth" -d "Google and GitHub integration" -p high
```

## 🎨 Status Output

GitLike status shows both Git and todo context:

```bash
gitlike status
# On branch feature-auth
# Todo status: 1 active, 2 completed
#
# Changes to be committed:
#   modified:   auth.js
#
# Untracked files:
#   login.html
```

## 🔗 Git Integration

GitLike works seamlessly with existing Git repositories:

```bash
# In existing Git repo
gitlike init  # Adds GitLike todo tracking
gitlike git init  # Enable full Git integration
gitlike git status  # View integration status
```

## 📝 Advanced Features

### Commit Message Templates
Commits automatically include todo context:

```bash
# With active todo
gitlike commit -m "Add validation"
# Creates: "Add validation (working on #1: Implement authentication)"
```

### Todo History
```bash
gitlike todo history 1
# Shows all commits related to todo #1
```

### Branch Management
```bash
gitlike branch  # List branches with todo counts
#   main (2 todos)
# * feature-auth (1 todo)
#   feature-forms (3 todos)
```

## 🛠️ Configuration

GitLike stores data in `~/.tododata/repository.json` with full Git integration support.

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `gitlike checkout -b feature-name`
3. Create todo: `gitlike todo create "Add new feature"`
4. Start working: `gitlike todo start <id>`
5. Commit changes: `gitlike commit -m "Implement feature"`
6. Mark complete: `gitlike todo done <id>`
7. Push branch: `gitlike push origin feature-name`

## 📄 License

MIT License - see LICENSE file for details.

---

**GitLike CLI** - Where Git meets task management for better development workflow! 🚀
