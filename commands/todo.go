package commands

import (
	"fmt"
	"gitlike/models"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var TodoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo related commands with Git-like syntax",
}

// GitLike syntax: gitlike todo create
var todoCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new todo task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, " ")
		description, _ := cmd.Flags().GetString("description")
		priority, _ := cmd.Flags().GetString("priority")

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		newTodo := models.Todo{
			ID:          repo.NextTodoID,
			Title:       title,
			Description: description,
			Status:      "pending",
			Priority:    priority,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			BranchName:  currentBranch.Name,
		}

		// Add todo to current branch
		for i := range repo.Branches {
			if repo.Branches[i].Name == currentBranch.Name {
				repo.Branches[i].Todos = append(repo.Branches[i].Todos, newTodo)
				break
			}
		}

		repo.NextTodoID++

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Created todo #%d: %s\n", newTodo.ID, newTodo.Title)
		if priority != "" {
			fmt.Printf("Priority: %s\n", priority)
		}
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}
	},
}

// GitLike syntax: gitlike todo done
var todoDoneCmd = &cobra.Command{
	Use:   "done [todo_id]",
	Short: "Mark a todo as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		todoIDStr := args[0]
		todoID, err := strconv.Atoi(todoIDStr)
		if err != nil {
			fmt.Printf("Invalid todo ID: %s\n", todoIDStr)
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		// Find and update the todo
		found := false
		for i := range repo.Branches {
			if repo.Branches[i].Name == currentBranch.Name {
				for j := range repo.Branches[i].Todos {
					if repo.Branches[i].Todos[j].ID == todoID {
						repo.Branches[i].Todos[j].Status = "completed"
						repo.Branches[i].Todos[j].UpdatedAt = time.Now()
						if repo.Branches[i].Todos[j].CompletedAt == nil {
							now := time.Now()
							repo.Branches[i].Todos[j].CompletedAt = &now
						}
						found = true
						break
					}
				}
				break
			}
		}

		if !found {
			fmt.Printf("Todo #%d not found in current branch\n", todoID)
			return
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("✅ Todo #%d marked as completed\n", todoID)
		
		// Suggest next steps
		fmt.Println("💡 Ready to commit your work:")
		fmt.Printf("   gitlike add .\n")
		fmt.Printf("   gitlike commit -m \"Complete todo #%d\"\n", todoID)
	},
}

var todoAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new todo (alias for create)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Just call todoCreateCmd
		todoCreateCmd.Run(cmd, args)
	},
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List todos in current branch",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		fmt.Printf("Todos in branch '%s':\n", currentBranch.Name)
		if len(currentBranch.Todos) == 0 {
			fmt.Println("  No todos found")
			return
		}

		for _, todo := range currentBranch.Todos {
			status := "⏳"
			switch todo.Status {
			case "completed":
				status = "✅"
			case "in-progress":
				status = "🔄"
			}

			activeIndicator := ""
			if todo.IsActive {
				activeIndicator = " 🎯"
			}

			commitInfo := ""
			if len(todo.Commits) > 0 {
				commitInfo = fmt.Sprintf(" (%d commits)", len(todo.Commits))
			}

			fmt.Printf("  %s #%d [%s] %s%s%s\n", status, todo.ID, todo.Priority, todo.Title, activeIndicator, commitInfo)
			if todo.Description != "" {
				fmt.Printf("    📝 %s\n", todo.Description)
			}
		}
	},
}

var todoUpdateCmd = &cobra.Command{
	Use:   "update [id] [status]",
	Short: "Update todo status (pending, in-progress, completed)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid todo ID: %s\n", args[0])
			return
		}

		status := args[1]
		if status != "pending" && status != "in-progress" && status != "completed" {
			fmt.Println("Status must be: pending, in-progress, or completed")
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		// Find and update todo
		found := false
		for i := range repo.Branches {
			if repo.Branches[i].Name == currentBranch.Name {
				for j := range repo.Branches[i].Todos {
					if repo.Branches[i].Todos[j].ID == id {
						repo.Branches[i].Todos[j].Status = status
						repo.Branches[i].Todos[j].UpdatedAt = time.Now()

						// Set completed timestamp
						if status == "completed" {
							now := time.Now()
							repo.Branches[i].Todos[j].CompletedAt = &now
							repo.Branches[i].Todos[j].IsActive = false
						}

						found = true
						break
					}
				}
				break
			}
		}

		if !found {
			fmt.Printf("Todo #%d not found in current branch\n", id)
			return
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Updated todo #%d status to: %s\n", id, status)
	},
}

var todoStartCmd = &cobra.Command{
	Use:   "start [id]",
	Short: "Start working on a todo (makes it the active task)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid todo ID: %s\n", args[0])
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		// First, deactivate all todos in current branch
		found := false
		for i := range repo.Branches {
			if repo.Branches[i].Name == currentBranch.Name {
				for j := range repo.Branches[i].Todos {
					repo.Branches[i].Todos[j].IsActive = false

					// Find and activate the target todo
					if repo.Branches[i].Todos[j].ID == id {
						repo.Branches[i].Todos[j].IsActive = true
						repo.Branches[i].Todos[j].Status = "in-progress"
						repo.Branches[i].Todos[j].UpdatedAt = time.Now()

						// Set started timestamp if not already set
						if repo.Branches[i].Todos[j].StartedAt == nil {
							now := time.Now()
							repo.Branches[i].Todos[j].StartedAt = &now
						}

						found = true
					}
				}
				break
			}
		}

		if !found {
			fmt.Printf("Todo #%d not found in current branch\n", id)
			return
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("🚀 Started working on todo #%d\n", id)
		fmt.Println("💡 Tip: All commits you make now will be automatically linked to this task!")
	},
}

var todoStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop working on the current active todo",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		// Find and deactivate active todo
		found := false
		var stoppedTodo *models.Todo
		for i := range repo.Branches {
			if repo.Branches[i].Name == currentBranch.Name {
				for j := range repo.Branches[i].Todos {
					if repo.Branches[i].Todos[j].IsActive {
						repo.Branches[i].Todos[j].IsActive = false
						repo.Branches[i].Todos[j].Status = "pending"
						repo.Branches[i].Todos[j].UpdatedAt = time.Now()
						stoppedTodo = &repo.Branches[i].Todos[j]
						found = true
						break
					}
				}
				break
			}
		}

		if !found {
			fmt.Println("No active todo found")
			return
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("⏸️  Stopped working on todo #%d: %s\n", stoppedTodo.ID, stoppedTodo.Title)
	},
}

var todoActiveCmd = &cobra.Command{
	Use:   "active",
	Short: "Show the currently active todo",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		// Find active todo
		var activeTodo *models.Todo
		for _, todo := range currentBranch.Todos {
			if todo.IsActive {
				activeTodo = &todo
				break
			}
		}

		if activeTodo == nil {
			fmt.Println("No active todo. Use 'gitlike todo start <id>' to start working on a task.")
			return
		}

		fmt.Printf("🎯 Currently working on:\n")
		fmt.Printf("  #%d [%s] %s\n", activeTodo.ID, activeTodo.Priority, activeTodo.Title)
		if activeTodo.Description != "" {
			fmt.Printf("  Description: %s\n", activeTodo.Description)
		}
		if activeTodo.StartedAt != nil {
			duration := time.Since(*activeTodo.StartedAt).Truncate(time.Minute)
			fmt.Printf("  Working for: %s\n", duration)
		}
		fmt.Printf("  Commits: %d\n", len(activeTodo.Commits))

		if len(activeTodo.Commits) > 0 {
			fmt.Println("  Recent commits:")
			// Show last few commits for this todo
			for i, commitID := range activeTodo.Commits {
				if i >= 3 { // Show max 3 recent commits
					break
				}
				for _, commit := range repo.Commits {
					if commit.ID == commitID {
						fmt.Printf("    - %s: %s\n", commitID[:8], commit.Message)
						break
					}
				}
			}
		}
	},
}

var todoHistoryCmd = &cobra.Command{
	Use:   "history [id]",
	Short: "Show commit history for a specific todo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid todo ID: %s\n", args[0])
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		currentBranch := storage_instance.GetCurrentBranch(repo)
		if currentBranch == nil {
			fmt.Println("No current branch found")
			return
		}

		// Find the todo
		var targetTodo *models.Todo
		for _, todo := range currentBranch.Todos {
			if todo.ID == id {
				targetTodo = &todo
				break
			}
		}

		if targetTodo == nil {
			fmt.Printf("Todo #%d not found in current branch\n", id)
			return
		}

		fmt.Printf("📋 Todo #%d: %s\n", targetTodo.ID, targetTodo.Title)
		fmt.Printf("Status: %s | Priority: %s\n", targetTodo.Status, targetTodo.Priority)

		if targetTodo.StartedAt != nil {
			fmt.Printf("Started: %s\n", targetTodo.StartedAt.Format("2006-01-02 15:04"))
		}

		if targetTodo.CompletedAt != nil {
			fmt.Printf("Completed: %s\n", targetTodo.CompletedAt.Format("2006-01-02 15:04"))
		}

		if len(targetTodo.Commits) == 0 {
			fmt.Println("\n📝 No commits linked to this todo yet")
			if targetTodo.IsActive {
				fmt.Println("💡 This is your active task - commits you make will be automatically linked!")
			}
			return
		}

		fmt.Printf("\n📝 Commit History (%d commits):\n", len(targetTodo.Commits))

		// Show commits in reverse order (newest first)
		for i := len(targetTodo.Commits) - 1; i >= 0; i-- {
			commitID := targetTodo.Commits[i]
			for _, commit := range repo.Commits {
				if commit.ID == commitID {
					fmt.Printf("  🔸 %s (%s)\n", commit.ID[:8], commit.CreatedAt.Format("Jan 2, 15:04"))
					fmt.Printf("     %s\n", commit.Message)
					fmt.Printf("     by %s\n", commit.Author)
					break
				}
			}
		}

		if len(targetTodo.Commits) > 0 {
			// Calculate development time from first commit
			var firstCommit *models.Commit
			for _, commit := range repo.Commits {
				if commit.ID == targetTodo.Commits[0] {
					firstCommit = &commit
					break
				}
			}
			if firstCommit != nil {
				duration := time.Since(firstCommit.CreatedAt).Truncate(time.Hour)
				fmt.Printf("\n⏱️  Development time: %s\n", duration)
			}
		}
	},
}

func init() {
	// Add flags for create and add commands
	todoCreateCmd.Flags().StringP("description", "d", "", "Todo description")
	todoCreateCmd.Flags().StringP("priority", "p", "medium", "Todo priority (low, medium, high)")
	
	todoAddCmd.Flags().StringP("description", "d", "", "Todo description")
	todoAddCmd.Flags().StringP("priority", "p", "medium", "Todo priority (low, medium, high)")

	// Add all subcommands to TodoCmd
	TodoCmd.AddCommand(todoCreateCmd)  // gitlike todo create
	TodoCmd.AddCommand(todoDoneCmd)    // gitlike todo done
	TodoCmd.AddCommand(todoAddCmd)     // gitlike todo add (alias)
	TodoCmd.AddCommand(todoListCmd)    // gitlike todo list
	TodoCmd.AddCommand(todoUpdateCmd)  // gitlike todo update
	TodoCmd.AddCommand(todoStartCmd)   // gitlike todo start
	TodoCmd.AddCommand(todoStopCmd)    // gitlike todo stop
	TodoCmd.AddCommand(todoActiveCmd)  // gitlike todo active
	TodoCmd.AddCommand(todoHistoryCmd) // gitlike todo history
}
