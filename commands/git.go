package commands

import (
	"fmt"
	"os/exec"
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
			
			// Get push status before pushing
			pushStatus, err := gitService.GetPushStatus()
			if err == nil {
				if unpushedCount, ok := pushStatus["unpushed_commits"].(int); ok && unpushedCount > 0 {
					fmt.Printf("📤 Pushing %d commits to remote...\n", unpushedCount)
				}
			}
			
			output, err := gitService.PushToRemoteWithDetails()
			if err != nil {
				fmt.Printf("Git push failed: %v\n", err)
				if strings.Contains(err.Error(), "no upstream branch") {
					fmt.Println("💡 Tip: Use --set-upstream to establish tracking")
				}
			} else {
				fmt.Println("✅ Pushed to Git remote")
				// Show push output if it contains useful info
				if strings.Contains(output, "->") {
					lines := strings.Split(output, "\n")
					for _, line := range lines {
						if strings.Contains(line, "->") || strings.Contains(line, "branch") {
							fmt.Printf("   %s\n", line)
						}
					}
				}
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

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push Git commits to remote repository",
	Run: func(cmd *cobra.Command, args []string) {
		if !gitService.IsGitRepo() {
			fmt.Println("Error: Not in a Git repository")
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		if !repo.GitIntegration.Enabled {
			fmt.Println("Git integration not enabled. Run 'todo git init' first.")
			return
		}

		// Get push status
		fmt.Println("Checking push status...")
		pushStatus, err := gitService.GetPushStatus()
		if err != nil {
			fmt.Printf("Error getting push status: %v\n", err)
			return
		}

		currentBranch := pushStatus["current_branch"].(string)
		fmt.Printf("📍 Current branch: %s\n", currentBranch)

		if remoteURL, ok := pushStatus["remote_url"].(string); ok {
			fmt.Printf("🌐 Remote: %s\n", remoteURL)
		}

		// Check for uncommitted changes
		if uncommittedCount, ok := pushStatus["uncommitted_changes"].(int); ok && uncommittedCount > 0 {
			fmt.Printf("⚠️  You have %d uncommitted changes\n", uncommittedCount)
			
			if changedFiles, ok := pushStatus["changed_files"].([]string); ok && len(changedFiles) <= 5 {
				for _, file := range changedFiles {
					fmt.Printf("   - %s\n", file)
				}
			}
			
			autoCommit, _ := cmd.Flags().GetBool("commit")
			if autoCommit {
				fmt.Println("🔄 Auto-committing changes...")
				commitMsg := "Auto-commit before push"
				err = gitService.CommitChanges(commitMsg)
				if err != nil {
					fmt.Printf("Auto-commit failed: %v\n", err)
					return
				}
				fmt.Println("✅ Changes committed")
			} else {
				fmt.Println("💡 Use --commit to auto-commit changes before push")
				fmt.Println("   Or manually commit with: git add . && git commit -m \"message\"")
			}
		}

		// Check for unpushed commits
		unpushedCount := 0
		if count, ok := pushStatus["unpushed_commits"].(int); ok {
			unpushedCount = count
		}

		if unpushedCount == 0 {
			fmt.Println("✨ Nothing to push - branch is up to date")
			return
		}

		fmt.Printf("📤 Found %d unpushed commits\n", unpushedCount)
		
		// Show commits to be pushed
		if commits, ok := pushStatus["commits"].([]models.Commit); ok && len(commits) > 0 {
			fmt.Println("Commits to push:")
			for i, commit := range commits {
				if i >= 3 { // Limit display to first 3 commits
					fmt.Printf("   ... and %d more commits\n", len(commits)-3)
					break
				}
				fmt.Printf("   %s %s (by %s)\n", commit.ID, commit.Message, commit.Author)
			}
		}

		// Perform the push
		fmt.Println("🚀 Pushing to remote...")
		output, err := gitService.PushToRemoteWithDetails()
		if err != nil {
			fmt.Printf("❌ Push failed: %v\n", err)
			
			// Provide helpful error messages
			errorStr := err.Error()
			if strings.Contains(errorStr, "no upstream branch") {
				fmt.Println("💡 No upstream branch set. Setting up tracking...")
				// This is handled in PushToRemoteWithDetails
			} else if strings.Contains(errorStr, "rejected") {
				fmt.Println("💡 Push rejected. Try pulling first: git pull")
			} else if strings.Contains(errorStr, "authentication") {
				fmt.Println("💡 Authentication failed. Check your Git credentials")
			}
			return
		}

		fmt.Println("✅ Successfully pushed to remote!")
		
		// Parse and display push output
		if output != "" {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				
				if strings.Contains(line, "->") {
					fmt.Printf("   📌 %s\n", line)
				} else if strings.Contains(line, "branch") && strings.Contains(line, "set up to track") {
					fmt.Printf("   🔗 %s\n", line)
				} else if strings.Contains(line, "objects") || strings.Contains(line, "Delta compression") {
					fmt.Printf("   📊 %s\n", line)
				}
			}
		}

		// Update last sync time
		repo.GitIntegration.LastGitSync = time.Now()
		storage_instance.SaveRepository(repo)
	},
}

var gitPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from Git remote repository",
	Run: func(cmd *cobra.Command, args []string) {
		if !gitService.IsGitRepo() {
			fmt.Println("Error: Not in a Git repository")
			return
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		if !repo.GitIntegration.Enabled {
			fmt.Println("Git integration not enabled. Run 'todo git init' first.")
			return
		}

		currentBranch, _ := gitService.GetCurrentBranch()
		fmt.Printf("📍 Pulling changes for branch: %s\n", currentBranch)

		if remoteURL, err := gitService.GetRemoteURL(); err == nil {
			fmt.Printf("🌐 From remote: %s\n", remoteURL)
		}

		// Check for uncommitted changes
		changedFiles, err := gitService.GetChangedFiles()
		if err == nil && len(changedFiles) > 0 {
			fmt.Printf("⚠️  You have %d uncommitted changes\n", len(changedFiles))
			fmt.Println("💡 Consider committing or stashing changes before pull")
			
			stash, _ := cmd.Flags().GetBool("stash")
			if stash {
				fmt.Println("📦 Stashing changes...")
				stashCmd := exec.Command("git", "stash", "push", "-m", "Auto-stash before pull")
				stashCmd.Dir = gitService.GetRepoPath()
				_, err := stashCmd.CombinedOutput()
				if err != nil {
					fmt.Printf("Stash failed: %v\n", err)
					return
				}
				fmt.Println("✅ Changes stashed")
			}
		}

		// Perform the pull
		fmt.Println("⬇️  Pulling from remote...")
		err = gitService.PullFromRemote()
		if err != nil {
			fmt.Printf("❌ Pull failed: %v\n", err)
			
			if strings.Contains(err.Error(), "merge conflict") {
				fmt.Println("💡 Merge conflicts detected. Resolve conflicts and commit.")
			} else if strings.Contains(err.Error(), "diverged") {
				fmt.Println("💡 Branches have diverged. Consider rebasing or merging.")
			}
			return
		}

		fmt.Println("✅ Successfully pulled from remote!")

		// Auto-sync todo branches if requested
		autoSync, _ := cmd.Flags().GetBool("sync")
		if autoSync {
			fmt.Println("🔄 Auto-syncing todo branches...")
			err = gitService.SyncWithGit(repo)
			if err != nil {
				fmt.Printf("Warning: Todo sync failed: %v\n", err)
			} else {
				fmt.Println("✅ Todo branches synchronized")
			}
		}

		// Update last sync time
		repo.GitIntegration.LastGitSync = time.Now()
		storage_instance.SaveRepository(repo)
	},
}

func init() {
	// Add flags
	gitCommitCmd.Flags().BoolP("push", "p", false, "Auto-push to Git remote after commit")
	gitPushCmd.Flags().BoolP("commit", "c", false, "Auto-commit changes before push")
	gitPullCmd.Flags().BoolP("stash", "s", false, "Auto-stash changes before pull")
	gitPullCmd.Flags().BoolP("sync", "y", false, "Auto-sync todo branches after pull")

	// Add subcommands
	GitCmd.AddCommand(gitInitCmd)
	GitCmd.AddCommand(gitStatusCmd)
	GitCmd.AddCommand(gitSyncCmd)
	GitCmd.AddCommand(gitCommitCmd)
	GitCmd.AddCommand(gitBranchCmd)
	GitCmd.AddCommand(gitPushCmd)
	GitCmd.AddCommand(gitPullCmd)
}
