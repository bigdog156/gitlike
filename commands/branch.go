package commands

import (
	"fmt"
	"gitlike/git"
	"gitlike/models"
	"gitlike/remote"
	"time"

	"github.com/spf13/cobra"
)

var BranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Branch related commands",
}

var branchCreateCmd = &cobra.Command{
	Use:   "create [branch_name]",
	Short: "Create a new branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Check if branch already exists
		if storage_instance.GetBranchByName(repo, branchName) != nil {
			fmt.Printf("Branch '%s' already exists\n", branchName)
			return
		}

		// Create new branch
		newBranch := models.Branch{
			Name:      branchName,
			CreatedAt: time.Now(),
			IsActive:  false,
			Todos:     []models.Todo{},
		}

		repo.Branches = append(repo.Branches, newBranch)

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Created branch: %s\n", branchName)
	},
}

var branchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all branches",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		fmt.Println("Branches:")
		for _, branch := range repo.Branches {
			current := ""
			if branch.Name == repo.CurrentBranch {
				current = " (current)"
			}
			fmt.Printf("  %s%s - %d todos\n", branch.Name, current, len(branch.Todos))
		}
	},
}

var branchSwitchCmd = &cobra.Command{
	Use:   "switch [branch_name]",
	Short: "Switch to a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		sync, _ := cmd.Flags().GetBool("sync")

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Check if branch exists locally
		if storage_instance.GetBranchByName(repo, branchName) == nil {
			// Try to pull from remote if sync is enabled
			if sync && len(repo.Remotes) > 0 {
				fmt.Printf("Branch '%s' not found locally, trying to sync from remote...\n", branchName)

				// Use first remote (usually origin)
				targetRemote := repo.Remotes[0]
				remoteService := remote.NewRemoteService()

				remoteRepo, err := remoteService.PullRepository(targetRemote)
				if err != nil {
					fmt.Printf("Failed to sync from remote: %v\n", err)
					return
				}

				// Check if branch exists in remote
				remoteBranch := storage_instance.GetBranchByName(remoteRepo, branchName)
				if remoteBranch != nil {
					// Add remote branch to local repo
					repo.Branches = append(repo.Branches, *remoteBranch)
					fmt.Printf("Pulled branch '%s' from remote\n", branchName)
				} else {
					fmt.Printf("Branch '%s' does not exist locally or on remote\n", branchName)
					return
				}
			} else {
				fmt.Printf("Branch '%s' does not exist\n", branchName)
				return
			}
		}

		repo.CurrentBranch = branchName

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Switched to branch: %s\n", branchName)

		// Auto-sync with Git if enabled
		if repo.GitIntegration.Enabled && repo.GitIntegration.AutoSync {
			gitService := git.NewGitService()
			if gitService.IsGitRepo() {
				fmt.Println("Auto-syncing with Git...")
				err = gitService.CheckoutBranch(branchName)
				if err != nil {
					fmt.Printf("Git checkout failed: %v\n", err)
				} else {
					fmt.Printf("✅ Git branch synchronized: %s\n", branchName)
				}
			}
		} // Auto-sync if requested
		if sync && len(repo.Remotes) > 0 {
			fmt.Println("Syncing with remote...")
			targetRemote := repo.Remotes[0]
			remoteService := remote.NewRemoteService()

			remoteRepo, err := remoteService.PullRepository(targetRemote)
			if err == nil {
				mergedRepo := remoteService.MergeRepositories(repo, remoteRepo)
				storage_instance.SaveRepository(mergedRepo)
				fmt.Println("Synced with remote")
			}
		}
	},
}

var branchDeleteCmd = &cobra.Command{
	Use:   "delete [branch_name]",
	Short: "Delete a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		force, _ := cmd.Flags().GetBool("force")

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Prevent deleting the current branch
		if repo.CurrentBranch == branchName {
			fmt.Printf("❌ Cannot delete the current branch '%s'. Switch to another branch first.\n", branchName)
			return
		}

		// Prevent deleting main branch unless forced
		if branchName == "main" && !force {
			fmt.Printf("❌ Cannot delete 'main' branch. Use --force flag if you really want to delete it.\n")
			return
		}

		// Find the branch to delete
		branchIndex := -1
		var branchToDelete *models.Branch
		for i, branch := range repo.Branches {
			if branch.Name == branchName {
				branchIndex = i
				branchToDelete = &branch
				break
			}
		}

		if branchIndex == -1 {
			fmt.Printf("❌ Branch '%s' does not exist\n", branchName)
			return
		}

		// Check if branch has todos and warn user
		if len(branchToDelete.Todos) > 0 && !force {
			fmt.Printf("⚠️  Branch '%s' has %d todos. Use --force flag to delete anyway.\n", branchName, len(branchToDelete.Todos))
			return
		}

		// Remove the branch from the slice
		repo.Branches = append(repo.Branches[:branchIndex], repo.Branches[branchIndex+1:]...)

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		if len(branchToDelete.Todos) > 0 {
			fmt.Printf("✅ Deleted branch '%s' (removed %d todos)\n", branchName, len(branchToDelete.Todos))
		} else {
			fmt.Printf("✅ Deleted branch '%s'\n", branchName)
		}
	},
}

func init() {
	// Add flags
	branchSwitchCmd.Flags().BoolP("sync", "s", false, "Sync with remote when switching branches")
	branchDeleteCmd.Flags().BoolP("force", "f", false, "Force delete branch even if it has todos or is the main branch")

	BranchCmd.AddCommand(branchCreateCmd)
	BranchCmd.AddCommand(branchListCmd)
	BranchCmd.AddCommand(branchSwitchCmd)
	BranchCmd.AddCommand(branchDeleteCmd)
}
