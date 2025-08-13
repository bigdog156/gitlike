package main

import (
	"fmt"
	"gitlike/commands"
	"os"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags during build)
var version = "dev-1.0.0"

// Root command
var rootCmd = &cobra.Command{
	Use:   "gitlike",
	Short: "GitLike CLI app with branch, commit, and merge functionality",
	Long: `A GitLike CLI application that helps developers track tasks with Git-like branch, commit, and merge operations.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`GitLike CLI - Git-like task management

Usage:
  gitlike [command]

Available Commands:
  branch      Branch related commands (create, list, switch, delete)
  todo        Todo related commands (add, list, update)
  commit      Commit related commands (create, list, show)
  merge       Merge a branch into current branch
  git         Git integration commands (push, status)
  remote      Remote repository commands (add, list, remove)
  help        Help about any command

Use "gitlike [command] --help" for more information about a command.
  
Quick Commands:
  gitlike commit "message"     - Create GitLike commit and Git commit
  gitlike todo start <id>      - Start working on a task
  gitlike todo active          - Show current active task`)
	},
}

// Version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GitLike CLI v%s\n", version)
	},
}

func main() {
	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Add all command groups
	rootCmd.AddCommand(commands.BranchCmd)
	rootCmd.AddCommand(commands.TodoCmd)
	rootCmd.AddCommand(commands.CommitCmd)
	// rootCmd.AddCommand(commands.MergeCmd)
	// rootCmd.AddCommand(commands.RemoteCmd)
	rootCmd.AddCommand(commands.GitCmd)

	// Add standalone remote commands
	// rootCmd.AddCommand(commands.PushCmd)
	// rootCmd.AddCommand(commands.PullCmd)
	// rootCmd.AddCommand(commands.FetchCmd)
	// rootCmd.AddCommand(commands.SyncCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
