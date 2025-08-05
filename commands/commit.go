package commands

import (
	"crypto/sha1"
	"fmt"
	"gitlike/models"
	"os/user"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit related commands",
}

var commitCreateCmd = &cobra.Command{
	Use:   "create [message]",
	Short: "Create a commit with completed todos",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")

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

		// Find completed todos
		var completedTodos []int
		for _, todo := range currentBranch.Todos {
			if todo.Status == "completed" {
				completedTodos = append(completedTodos, todo.ID)
			}
		}

		if len(completedTodos) == 0 {
			fmt.Println("No completed todos to commit")
			return
		}

		// Get current user
		currentUser, _ := user.Current()
		author := currentUser.Username
		if author == "" {
			author = "unknown"
		}

		// Generate commit ID
		hash := sha1.New()
		hash.Write([]byte(fmt.Sprintf("%s-%s-%d", message, currentBranch.Name, time.Now().UnixNano())))
		commitID := fmt.Sprintf("%x", hash.Sum(nil))[:8]

		// Create commit
		commit := models.Commit{
			ID:        commitID,
			Message:   message,
			Branch:    currentBranch.Name,
			Todos:     completedTodos,
			CreatedAt: time.Now(),
			Author:    author,
		}

		repo.Commits = append(repo.Commits, commit)

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Created commit %s: %s\n", commitID, message)
		fmt.Printf("Committed %d completed todos\n", len(completedTodos))
	},
}

var commitListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all commits",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		if len(repo.Commits) == 0 {
			fmt.Println("No commits found")
			return
		}

		fmt.Println("Commits:")
		for i := len(repo.Commits) - 1; i >= 0; i-- {
			commit := repo.Commits[i]
			fmt.Printf("  %s [%s] %s by %s (%d todos)\n",
				commit.ID, commit.Branch, commit.Message, commit.Author, len(commit.Todos))
			fmt.Printf("    %s\n", commit.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

var commitShowCmd = &cobra.Command{
	Use:   "show [commit_id]",
	Short: "Show commit details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitID := args[0]

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Find commit
		var commit *models.Commit
		for i := range repo.Commits {
			if repo.Commits[i].ID == commitID {
				commit = &repo.Commits[i]
				break
			}
		}

		if commit == nil {
			fmt.Printf("Commit %s not found\n", commitID)
			return
		}

		fmt.Printf("Commit: %s\n", commit.ID)
		fmt.Printf("Message: %s\n", commit.Message)
		fmt.Printf("Branch: %s\n", commit.Branch)
		fmt.Printf("Author: %s\n", commit.Author)
		fmt.Printf("Date: %s\n", commit.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Todos (%d):\n", len(commit.Todos))

		// Find branch to get todo details
		branch := storage_instance.GetBranchByName(repo, commit.Branch)
		if branch != nil {
			for _, todoID := range commit.Todos {
				for _, todo := range branch.Todos {
					if todo.ID == todoID {
						fmt.Printf("  #%d %s - %s\n", todo.ID, todo.Title, todo.Description)
						break
					}
				}
			}
		}
	},
}

func init() {
	CommitCmd.AddCommand(commitCreateCmd)
	CommitCmd.AddCommand(commitListCmd)
	CommitCmd.AddCommand(commitShowCmd)
}
