package commands

import (
	"fmt"
	"strings"
	"time"
	"todo-cli/git"
	"todo-cli/models"

	"github.com/spf13/cobra"
)

var gitService = git.NewGitService()

var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git integration commands",
}

var gitInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Git integration",
	Run: func(cmd *cobra.Command, args []string) {
		if !gitService.IsGitRepo() {
			fmt.Println("Error: Not in a Git repository")
			fmt.Println("Please run this command from within a Git repository")
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Enable Git integration
		repo.GitIntegration.Enabled = true
		repo.GitIntegration.AutoSync = true
		repo.GitIntegration.RepoPath = gitService.GetRepoPath()

		// Get Git remote URL
		if remoteURL, err := gitService.GetRemoteURL(); err == nil {
			repo.GitIntegration.RemoteURL = remoteURL
		}

		// Set default commit template
		repo.GitIntegration.CommitTemplate = "todo: {{.Message}}"

		// Sync with Git
		err = gitService.SyncWithGit(repo)
		if err != nil {
			fmt.Printf("Warning: Could not sync with Git: %v\n", err)
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Println("✅ Git integration initialized successfully!")
		fmt.Printf("📁 Repository: %s\n", repo.GitIntegration.RepoPath)
		if repo.GitIntegration.RemoteURL != "" {
			fmt.Printf("🌐 Remote: %s\n", repo.GitIntegration.RemoteURL)
		}

		currentBranch, _ := gitService.GetCurrentBranch()
		fmt.Printf("🌿 Current branch: %s\n", currentBranch)
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Git integration status",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		fmt.Println("Git Integration Status:")
		fmt.Printf("Enabled: %v\n", repo.GitIntegration.Enabled)

		if !repo.GitIntegration.Enabled {
			fmt.Println("Run 'todo git init' to enable Git integration")
			return
		}

		fmt.Printf("Auto-sync: %v\n", repo.GitIntegration.AutoSync)
		fmt.Printf("Auto-commit: %v\n", repo.GitIntegration.AutoCommit)
		fmt.Printf("Repository path: %s\n", repo.GitIntegration.RepoPath)

		if repo.GitIntegration.RemoteURL != "" {
			fmt.Printf("Remote URL: %s\n", repo.GitIntegration.RemoteURL)
		}

		if gitService.IsGitRepo() {
			currentBranch, _ := gitService.GetCurrentBranch()
			fmt.Printf("Current Git branch: %s\n", currentBranch)

			changedFiles, err := gitService.GetChangedFiles()
			if err == nil {
				fmt.Printf("Changed files: %d\n", len(changedFiles))
				if len(changedFiles) > 0 && len(changedFiles) <= 5 {
					for _, file := range changedFiles {
						fmt.Printf("  - %s\n", file)
					}
				}
			}
		} else {
			fmt.Println("⚠️  Not currently in a Git repository")
		}

		if !repo.GitIntegration.LastGitSync.IsZero() {
			fmt.Printf("Last Git sync: %s\n", repo.GitIntegration.LastGitSync.Format("2006-01-02 15:04:05"))
		}
	},
}

var gitSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize with Git repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		if !repo.GitIntegration.Enabled {
			fmt.Println("Git integration not enabled. Run 'todo git init' first.")
			return
		}

		if !gitService.IsGitRepo() {
			fmt.Println("Error: Not in a Git repository")
			return
		}

		fmt.Println("Synchronizing with Git...")

		err = gitService.SyncWithGit(repo)
		if err != nil {
			fmt.Printf("Error syncing with Git: %v\n", err)
			return
		}

		// Update last sync time
		repo.GitIntegration.LastGitSync = time.Now()

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		currentBranch, _ := gitService.GetCurrentBranch()
		fmt.Printf("✅ Synchronized with Git branch: %s\n", currentBranch)
		fmt.Printf("📊 Todo branches: %d\n", len(repo.Branches))
	},
}

var gitCommitCmd = &cobra.Command{
	Use:   "commit [message]",
	Short: "Commit completed todos to Git",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		if !repo.GitIntegration.Enabled {
			fmt.Println("Git integration not enabled. Run 'todo git init' first.")
			return
		}

		if !gitService.IsGitRepo() {
			fmt.Println("Error: Not in a Git repository")
			return
		}

		// First create a todo commit
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

		// Create todo commit first
		commitCreateCmd.Run(cmd, args)

		// Then create Git commit
		fmt.Println("Creating Git commit...")

		gitMessage := fmt.Sprintf("todo: %s\n\nCompleted todos: %v", message, completedTodos)
		err = gitService.CommitChanges(gitMessage)
		if err != nil {
			fmt.Printf("Git commit failed: %v\n", err)
			return
		}

		fmt.Printf("✅ Created Git commit with %d completed todos\n", len(completedTodos))

		// Auto-push if requested
		autoPush, _ := cmd.Flags().GetBool("push")
		if autoPush {
			fmt.Println("Pushing to Git remote...")
			err = gitService.PushToRemote()
			if err != nil {
				fmt.Printf("Git push failed: %v\n", err)
			} else {
				fmt.Println("✅ Pushed to Git remote")
			}
		}
	},
}

var gitBranchCmd = &cobra.Command{
	Use:   "branch [action] [name]",
	Short: "Git branch operations synchronized with todos",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]

		if !gitService.IsGitRepo() {
			fmt.Println("Error: Not in a Git repository")
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		switch action {
		case "create":
			if len(args) < 2 {
				fmt.Println("Usage: todo git branch create [branch_name]")
				return
			}
			branchName := args[1]

			// Create Git branch
			err = gitService.CreateBranch(branchName)
			if err != nil {
				fmt.Printf("Error creating Git branch: %v\n", err)
				return
			}

			// Create todo branch
			newBranch := models.Branch{
				Name:      branchName,
				CreatedAt: time.Now(),
				IsActive:  true,
				Todos:     []models.Todo{},
			}
			repo.Branches = append(repo.Branches, newBranch)
			repo.CurrentBranch = branchName

			err = storage_instance.SaveRepository(repo)
			if err != nil {
				fmt.Printf("Error saving repository: %v\n", err)
				return
			}

			fmt.Printf("✅ Created and switched to branch: %s (Git + Todo)\n", branchName)

		case "switch", "checkout":
			if len(args) < 2 {
				fmt.Println("Usage: todo git branch switch [branch_name]")
				return
			}
			branchName := args[1]

			// Switch Git branch
			err = gitService.CheckoutBranch(branchName)
			if err != nil {
				fmt.Printf("Error switching Git branch: %v\n", err)
				return
			}

			// Switch todo branch
			if storage_instance.GetBranchByName(repo, branchName) == nil {
				// Create todo branch if it doesn't exist
				newBranch := models.Branch{
					Name:      branchName,
					CreatedAt: time.Now(),
					IsActive:  false,
					Todos:     []models.Todo{},
				}
				repo.Branches = append(repo.Branches, newBranch)
			}

			repo.CurrentBranch = branchName

			err = storage_instance.SaveRepository(repo)
			if err != nil {
				fmt.Printf("Error saving repository: %v\n", err)
				return
			}

			fmt.Printf("✅ Switched to branch: %s (Git + Todo)\n", branchName)

		case "list":
			gitBranches, err := gitService.GetAllBranches()
			if err != nil {
				fmt.Printf("Error getting Git branches: %v\n", err)
				return
			}

			currentBranch, _ := gitService.GetCurrentBranch()

			fmt.Println("Branches (Git + Todo sync):")
			for _, branch := range gitBranches {
				current := ""
				if branch == currentBranch {
					current = " (current)"
				}

				// Get todo count for this branch
				todoCount := 0
				todoBranch := storage_instance.GetBranchByName(repo, branch)
				if todoBranch != nil {
					todoCount = len(todoBranch.Todos)
				}

				fmt.Printf("  %s%s - %d todos\n", branch, current, todoCount)
			}

		default:
			fmt.Printf("Unknown action: %s\n", action)
			fmt.Println("Available actions: create, switch, checkout, list")
		}
	},
}

func init() {
	// Add flags
	gitCommitCmd.Flags().BoolP("push", "p", false, "Auto-push to Git remote after commit")

	// Add subcommands
	GitCmd.AddCommand(gitInitCmd)
	GitCmd.AddCommand(gitStatusCmd)
	GitCmd.AddCommand(gitSyncCmd)
	GitCmd.AddCommand(gitCommitCmd)
	GitCmd.AddCommand(gitBranchCmd)
}
