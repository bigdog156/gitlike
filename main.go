package main

import (
	"fmt"
	"os"
	"todo-cli/commands"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo CLI app with branch, commit, and merge functionality",
	Long: `A todo CLI application that helps developers track tasks with Git-like branch, commit, and merge operations.
    
Examples:
  todo branch create feature-auth
  todo todo add "Implement user login" -d "Add JWT authentication" -p high
  todo todo update 1 completed
  todo commit create "Implement user authentication"
  todo merge feature-auth`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Todo CLI - Git-like task management

Usage:
  todo [command]

Available Commands:
  branch      Branch related commands (create, list, switch)
  todo        Todo related commands (add, list, update)
  commit      Commit related commands (create, list, show)
  merge       Merge a branch into current branch
  help        Help about any command

Use "todo [command] --help" for more information about a command.`)
	},
}

func main() {
	// Add all command groups
	rootCmd.AddCommand(commands.BranchCmd)
	rootCmd.AddCommand(commands.TodoCmd)
	rootCmd.AddCommand(commands.CommitCmd)
	rootCmd.AddCommand(commands.MergeCmd)
	rootCmd.AddCommand(commands.RemoteCmd)
	rootCmd.AddCommand(commands.GitCmd)

	// Add standalone remote commands
	rootCmd.AddCommand(commands.PushCmd)
	rootCmd.AddCommand(commands.PullCmd)
	rootCmd.AddCommand(commands.FetchCmd)
	rootCmd.AddCommand(commands.SyncCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
