#!/bin/bash

# GitLike CLI - Complete Git-like Workflow with Todo Tracking Demo

echo "🚀 GitLike CLI - Git-like Task Management Demo"
echo "============================================="
echo ""

# Build the application
cd /Users/lethachlam/Developer/todocli
go build -o gitlike .

echo "1. Initialize a new GitLike repository"
cd /tmp && rm -rf gitlike-demo && mkdir gitlike-demo && cd gitlike-demo
/Users/lethachlam/Developer/todocli/gitlike init
echo ""

echo "2. Create todos with Git-like syntax"
/Users/lethachlam/Developer/todocli/gitlike todo create "Implement user authentication" -d "Add JWT auth system" -p high
/Users/lethachlam/Developer/todocli/gitlike todo create "Create login form" -d "HTML form with validation" -p medium
/Users/lethachlam/Developer/todocli/gitlike todo create "Add password reset" -d "Email-based reset flow" -p low
echo ""

echo "3. List all todos"
/Users/lethachlam/Developer/todocli/gitlike todo list
echo ""

echo "4. Start working on a task (makes it active)"
/Users/lethachlam/Developer/todocli/gitlike todo start 1
echo ""

echo "5. Create some files and check status"
echo "# Auth System" > auth.js
echo "# Login Form" > login.html
echo "# Main App" > app.js
/Users/lethachlam/Developer/todocli/gitlike status
echo ""

echo "6. Add files to staging (Git-like syntax)"
/Users/lethachlam/Developer/todocli/gitlike add auth.js login.html
echo ""

echo "7. Check status after adding"
/Users/lethachlam/Developer/todocli/gitlike status
echo ""

echo "8. Commit with todo tracking (Git-like syntax)"
/Users/lethachlam/Developer/todocli/gitlike commit -m "Add authentication scaffolding"
echo ""

echo "9. Mark todo as completed"
/Users/lethachlam/Developer/todocli/gitlike todo done 1
echo ""

echo "10. Add remaining file and commit completion"
/Users/lethachlam/Developer/todocli/gitlike add app.js
/Users/lethachlam/Developer/todocli/gitlike commit -m "Complete user authentication"
echo ""

echo "11. Create and switch to feature branch (Git-like syntax)"
/Users/lethachlam/Developer/todocli/gitlike checkout -b feature-forms
echo ""

echo "12. Start working on next todo"
/Users/lethachlam/Developer/todocli/gitlike todo start 2
echo ""

echo "13. Create branch-specific work"
echo "<form>Login form</form>" > form.html
/Users/lethachlam/Developer/todocli/gitlike add form.html
/Users/lethachlam/Developer/todocli/gitlike commit -m "Add login form HTML"
echo ""

echo "14. View commit history with todo tracking"
/Users/lethachlam/Developer/todocli/gitlike log
echo ""

echo "15. Switch back to main branch"
/Users/lethachlam/Developer/todocli/gitlike checkout main
echo ""

echo "16. Show branch-specific todos"
echo "Todos in main branch:"
/Users/lethachlam/Developer/todocli/gitlike todo list
echo ""

echo "17. Switch to feature branch and show its todos"
/Users/lethachlam/Developer/todocli/gitlike checkout feature-forms
echo "Todos in feature-forms branch:"
/Users/lethachlam/Developer/todocli/gitlike todo list
echo ""

echo "✅ GitLike CLI Demo Complete!"
echo ""
echo "🎯 Key Features Demonstrated:"
echo "  • Git-identical syntax (init, add, commit, status, log, checkout)"  
echo "  • Automatic todo tracking with commits"
echo "  • Branch-specific todo management"
echo "  • Active task workflow (start/stop)"
echo "  • Todo completion tracking"
echo "  • Rich commit messages with todo context"
echo ""
echo "📝 Git-like Commands Available:"
echo "  gitlike init                    # Initialize repository"
echo "  gitlike add <files>             # Stage changes"
echo "  gitlike commit -m \"message\"     # Commit with todo tracking"
echo "  gitlike status                  # Show status with todo context"
echo "  gitlike log                     # Show commits with todo history"
echo "  gitlike checkout <branch>       # Switch branches"  
echo "  gitlike checkout -b <branch>    # Create and switch branch"
echo "  gitlike todo create \"title\"     # Create new todo"
echo "  gitlike todo done <id>          # Mark todo complete"
echo "  gitlike todo start <id>         # Start working on todo"
echo "  gitlike todo list               # List todos in current branch"
echo ""
