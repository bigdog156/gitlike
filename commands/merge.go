package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var MergeCmd = &cobra.Command{
	Use:   "merge [source_branch]",
	Short: "Merge a branch into current branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourceBranch := args[0]

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

		if currentBranch.Name == sourceBranch {
			fmt.Println("Cannot merge branch into itself")
			return
		}

		// Find source branch
		sourceB := storage_instance.GetBranchByName(repo, sourceBranch)
		if sourceB == nil {
			fmt.Printf("Branch '%s' does not exist\n", sourceBranch)
			return
		}

		// Merge todos from source branch
		mergedCount := 0
		for _, sourceTodo := range sourceB.Todos {
			// Check if todo already exists in current branch (by ID)
			exists := false
			for _, currentTodo := range currentBranch.Todos {
				if currentTodo.ID == sourceTodo.ID {
					exists = true
					break
				}
			}

			if !exists {
				// Update branch name and add to current branch
				mergedTodo := sourceTodo
				mergedTodo.BranchName = currentBranch.Name
				mergedTodo.UpdatedAt = time.Now()

				// Add to current branch
				for i := range repo.Branches {
					if repo.Branches[i].Name == currentBranch.Name {
						repo.Branches[i].Todos = append(repo.Branches[i].Todos, mergedTodo)
						mergedCount++
						break
					}
				}
			}
		}

		// Merge commits from source branch
		mergedCommits := 0
		for _, commit := range repo.Commits {
			if commit.Branch == sourceBranch {
				// Create a new commit entry for the current branch
				mergedCommit := commit
				mergedCommit.Branch = currentBranch.Name
				mergedCommit.Message = fmt.Sprintf("[MERGED from %s] %s", sourceBranch, commit.Message)
				mergedCommit.CreatedAt = time.Now()

				repo.Commits = append(repo.Commits, mergedCommit)
				mergedCommits++
			}
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Merged branch '%s' into '%s'\n", sourceBranch, currentBranch.Name)
		fmt.Printf("- %d todos merged\n", mergedCount)
		fmt.Printf("- %d commits merged\n", mergedCommits)

		// Ask if user wants to delete the source branch
		fmt.Printf("\nDelete source branch '%s'? (y/N): ", sourceBranch)
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			// Remove source branch
			for i, branch := range repo.Branches {
				if branch.Name == sourceBranch {
					repo.Branches = append(repo.Branches[:i], repo.Branches[i+1:]...)
					break
				}
			}

			err = storage_instance.SaveRepository(repo)
			if err != nil {
				fmt.Printf("Error deleting branch: %v\n", err)
				return
			}

			fmt.Printf("Deleted branch '%s'\n", sourceBranch)
		}
	},
}
