#!/bin/bash

# GitLike CLI - Simple Usage Examples

echo "GitLike CLI - Simple Examples"
echo "============================="

# Example 1: Basic workflow
echo "Example 1: Basic Task-Driven Development"
echo "----------------------------------------"
echo "gitlike init                                    # Initialize repository"
echo "gitlike todo create 'Add user auth' -p high    # Create high-priority todo"
echo "gitlike todo start 1                           # Start working on todo #1"
echo "echo 'console.log(\"auth\")' > auth.js         # Create file"
echo "gitlike add auth.js                            # Stage file"
echo "gitlike commit -m 'Add authentication'         # Commit with todo tracking"
echo "gitlike todo done 1                            # Mark todo complete"
echo ""

# Example 2: Branch workflow  
echo "Example 2: Feature Branch with Todos"
echo "------------------------------------"
echo "gitlike checkout -b feature-login               # Create feature branch"
echo "gitlike todo create 'Build login form'         # Branch-specific todo"
echo "gitlike todo start 2                           # Start working"
echo "echo '<form>Login</form>' > login.html         # Create form"
echo "gitlike add login.html                         # Stage changes"
echo "gitlike commit -m 'Add login form'             # Commit"
echo "gitlike checkout main                          # Switch to main"
echo "gitlike todo list                              # See main branch todos"
echo ""

# Example 3: Status and history
echo "Example 3: Viewing Status and History"
echo "-------------------------------------"
echo "gitlike status                                  # Git status + todo context"
echo "gitlike log                                     # Commit history + todos"
echo "gitlike todo active                            # Show current active todo"
echo "gitlike todo history 1                         # Show todo's commit history"
echo ""

# All available commands
echo "All GitLike Commands:"
echo "===================="
echo ""
echo "🎯 Core Git Commands (identical to Git):"
echo "  init, add, commit, status, log, checkout, push, pull"
echo ""
echo "📝 Todo Management:"
echo "  gitlike todo create \"title\" [-d desc] [-p priority]"
echo "  gitlike todo done <id>"
echo "  gitlike todo start <id>"
echo "  gitlike todo stop"
echo "  gitlike todo list"
echo "  gitlike todo active"
echo ""
echo "🌿 Branch Operations:"
echo "  gitlike branch                 # List branches"
echo "  gitlike checkout -b <branch>   # Create branch"
echo "  gitlike checkout <branch>      # Switch branch"
echo ""
