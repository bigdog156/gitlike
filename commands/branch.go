package commands

import (
	"fmt"
	"time"
	"todo-cli/models"
	"todo-cli/remote"

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

		// Auto-sync if requested
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

func init() {
	// Add flags
	branchSwitchCmd.Flags().BoolP("sync", "s", false, "Sync with remote when switching branches")

	BranchCmd.AddCommand(branchCreateCmd)
	BranchCmd.AddCommand(branchListCmd)
	BranchCmd.AddCommand(branchSwitchCmd)
}
