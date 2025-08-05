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
	Short: "Todo related commands",
}

var todoAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new todo",
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

		fmt.Printf("Added todo #%d: %s\n", newTodo.ID, newTodo.Title)
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
			fmt.Printf("  %s #%d [%s] %s - %s\n", status, todo.ID, todo.Priority, todo.Title, todo.Description)
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

func init() {
	todoAddCmd.Flags().StringP("description", "d", "", "Todo description")
	todoAddCmd.Flags().StringP("priority", "p", "medium", "Todo priority (low, medium, high)")

	TodoCmd.AddCommand(todoAddCmd)
	TodoCmd.AddCommand(todoListCmd)
	TodoCmd.AddCommand(todoUpdateCmd)
}
