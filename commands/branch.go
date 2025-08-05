package commands

import (
    "fmt"
    "time"
    "github.com/spf13/cobra"
    "todo-cli/models"
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
        
        repo, err := storage_instance.LoadRepository()
        if err != nil {
            fmt.Printf("Error loading repository: %v\n", err)
            return
        }
        
        // Check if branch exists
        if storage_instance.GetBranchByName(repo, branchName) == nil {
            fmt.Printf("Branch '%s' does not exist\n", branchName)
            return
        }
        
        repo.CurrentBranch = branchName
        
        err = storage_instance.SaveRepository(repo)
        if err != nil {
            fmt.Printf("Error saving repository: %v\n", err)
            return
        }
        
        fmt.Printf("Switched to branch: %s\n", branchName)
    },
}

func init() {
    BranchCmd.AddCommand(branchCreateCmd)
    BranchCmd.AddCommand(branchListCmd)
    BranchCmd.AddCommand(branchSwitchCmd)
}