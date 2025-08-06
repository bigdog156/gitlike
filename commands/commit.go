package commands

import (
	"crypto/sha1"
	"fmt"
	"gitlike/models"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Helper function to check if current directory is a Git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit related commands",
}

var commitCreateCmd = &cobra.Command{
	Use:   "create [message]",
	Short: "Create a commit with todos and optionally commit to Git",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")
		gitCommit, _ := cmd.Flags().GetBool("git")
		autoAdd, _ := cmd.Flags().GetBool("add")

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

		// Find completed todos and active todo
		var completedTodos []int
		var activeTodo *models.Todo

		for _, todo := range currentBranch.Todos {
			if todo.Status == "completed" {
				completedTodos = append(completedTodos, todo.ID)
			}
			if todo.IsActive {
				activeTodo = &todo
			}
		}

		// Determine which todos to include in commit
		var todosToInclude []int
		
		if activeTodo != nil {
			// Prioritize active todo
			todosToInclude = append(todosToInclude, activeTodo.ID)
			fmt.Printf("Including active todo #%d in commit: %s\n", activeTodo.ID, activeTodo.Title)
			
			// If there are also completed todos, ask if user wants to include them too
			if len(completedTodos) > 0 {
				fmt.Printf("Also found %d completed todos: %v\n", len(completedTodos), completedTodos)
				fmt.Print("Include completed todos in this commit too? (y/N): ")
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
					todosToInclude = append(todosToInclude, completedTodos...)
				}
			}
		} else if len(completedTodos) > 0 {
			// No active todo, use completed todos
			todosToInclude = completedTodos
			fmt.Printf("Including %d completed todos in commit\n", len(completedTodos))
		} else {
			fmt.Println("No todos to commit. Start working on a todo with 'gitlike todo start <id>' or complete one first.")
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
			Todos:     todosToInclude,
			CreatedAt: time.Now(),
			Author:    author,
		}

		// Set active todo if there is one
		if activeTodo != nil {
			commit.ActiveTodo = &activeTodo.ID
		}

		repo.Commits = append(repo.Commits, commit)

		// Link this commit to the todos (add commit ID to todo's commits list)
		for i := range repo.Branches {
			if repo.Branches[i].Name == currentBranch.Name {
				for j := range repo.Branches[i].Todos {
					for _, todoID := range todosToInclude {
						if repo.Branches[i].Todos[j].ID == todoID {
							// Add commit ID to todo's commit list if not already present
							found := false
							for _, commitID := range repo.Branches[i].Todos[j].Commits {
								if commitID == commit.ID {
									found = true
									break
								}
							}
							if !found {
								repo.Branches[i].Todos[j].Commits = append(repo.Branches[i].Todos[j].Commits, commit.ID)
							}
							break
						}
					}
				}
				break
			}
		}

		repo.Commits = append(repo.Commits, commit)

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("✅ Created commit %s: %s\n", commitID, message)
		fmt.Printf("📋 Linked to %d todos: %v\n", len(todosToInclude), todosToInclude)

		if activeTodo != nil {
			fmt.Printf("🎯 Active task: #%d %s\n", activeTodo.ID, activeTodo.Title)
		}

		// Git integration
		if gitCommit {
			fmt.Println("\n🔄 Committing to Git...")
			
			// Check if we're in a Git repository
			if !isGitRepo() {
				fmt.Println("⚠️  Not in a Git repository. Skipping Git commit.")
				return
			}

			// Auto-add files if requested
			if autoAdd {
				fmt.Println("📁 Adding all changes to Git...")
				gitAddCmd := exec.Command("git", "add", ".")
				gitAddCmd.Dir, _ = os.Getwd()
				if err := gitAddCmd.Run(); err != nil {
					fmt.Printf("⚠️  Failed to add files to Git: %v\n", err)
				} else {
					fmt.Println("✅ Added all changes to Git")
				}
			}

			// Create enhanced Git commit message with todo information
			var gitMessage strings.Builder
			gitMessage.WriteString(message)
			
			if len(todosToInclude) > 0 {
				gitMessage.WriteString("\n\n")
				if activeTodo != nil {
					gitMessage.WriteString(fmt.Sprintf("Active Task: #%d %s\n", activeTodo.ID, activeTodo.Title))
				}
				
				gitMessage.WriteString("GitLike Todos:\n")
				currentBranch := storage_instance.GetCurrentBranch(repo)
				for _, todoID := range todosToInclude {
					for _, todo := range currentBranch.Todos {
						if todo.ID == todoID {
							gitMessage.WriteString(fmt.Sprintf("- #%d [%s] %s", todo.ID, todo.Priority, todo.Title))
							if todo.Status == "completed" {
								gitMessage.WriteString(" ✅")
							} else if todo.IsActive {
								gitMessage.WriteString(" 🎯")
							}
							gitMessage.WriteString("\n")
							break
						}
					}
				}
				gitMessage.WriteString(fmt.Sprintf("\nGitLike Commit ID: %s", commitID))
			}

			// Execute Git commit
			gitCommitCmd := exec.Command("git", "commit", "-m", gitMessage.String())
			gitCommitCmd.Dir, _ = os.Getwd()
			gitCommitCmd.Stdout = os.Stdout
			gitCommitCmd.Stderr = os.Stderr
			
			if err := gitCommitCmd.Run(); err != nil {
				fmt.Printf("⚠️  Git commit failed: %v\n", err)
				fmt.Println("💡 You may need to stage your changes first with 'git add' or use the --add flag")
			} else {
				fmt.Println("✅ Successfully committed to Git!")
				
				// Get the actual Git commit hash
				gitHashCmd := exec.Command("git", "rev-parse", "HEAD")
				gitHashCmd.Dir, _ = os.Getwd()
				if gitHashOutput, err := gitHashCmd.Output(); err == nil {
					gitHash := strings.TrimSpace(string(gitHashOutput))
					fmt.Printf("🔗 Git commit: %s\n", gitHash[:8])
					
					// Store the Git hash in the GitLike commit for reference
					for i := range repo.Commits {
						if repo.Commits[i].ID == commitID {
							// We could add a GitHash field to the Commit model if needed
							break
						}
					}
				}
			}
		}
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
	// Add flags for Git integration
	commitCreateCmd.Flags().BoolP("git", "g", false, "Also commit to Git repository")
	commitCreateCmd.Flags().BoolP("add", "a", false, "Automatically add all changes before Git commit (requires --git)")
	
	CommitCmd.AddCommand(commitCreateCmd)
	CommitCmd.AddCommand(commitListCmd)
	CommitCmd.AddCommand(commitShowCmd)
}
