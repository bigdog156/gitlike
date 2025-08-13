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
	Short: "GitLike CLI - Git-like task management with todo tracking",
	Long: `A GitLike CLI application that mirrors Git syntax with integrated todo tracking.

GitLike provides the same commands as Git but with enhanced task management:
  - Track todos with commits
  - Link tasks to development work
  - Maintain Git-like workflow with todo context
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`GitLike CLI v` + version + ` - Git-like task management

Usage:
  gitlike <command> [<args>]

Git-like Commands:
  init            Create an empty GitLike repository
  clone <url>     Clone a repository
  add <file>      Add file contents to the index
  commit          Record changes with todo tracking
  push            Update remote refs with todo sync
  pull            Fetch and integrate with another repository
  fetch           Download objects and refs from another repository
  merge           Join development histories together
  branch          List, create, or delete branches
  checkout        Switch branches or restore working tree files
  status          Show the working tree status
  log             Show commit logs with todo history
  diff            Show changes between commits, commit and working tree, etc
  remote          Manage set of tracked repositories
  
Todo Commands:
  todo create     Create a new todo task
  todo done       Mark todo as completed
  todo list       List todos in current branch
  todo start      Start working on a todo
  todo stop       Stop working on current todo
  
Examples:
  gitlike init
  gitlike todo create "Implement user authentication"
  gitlike add .
  gitlike commit -m "Add login form"
  gitlike push origin main
  
Use "gitlike <command> --help" for more information on a specific command.`)
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

	// Add Git-like commands (primary interface)
	rootCmd.AddCommand(commands.InitCmd)         // gitlike init
	rootCmd.AddCommand(commands.AddCmd)          // gitlike add
	rootCmd.AddCommand(commands.CommitDirectCmd) // gitlike commit -m
	rootCmd.AddCommand(commands.StatusCmd)       // gitlike status
	rootCmd.AddCommand(commands.LogCmd)          // gitlike log
	rootCmd.AddCommand(commands.CheckoutCmd)     // gitlike checkout
	rootCmd.AddCommand(commands.PushCmd)         // gitlike push (from remote.go)
	rootCmd.AddCommand(commands.PullCmd)         // gitlike pull (from remote.go)

	// Add branch command (Git-like: gitlike branch)
	rootCmd.AddCommand(commands.BranchCmd)

	// Add todo command with Git-like subcommands
	rootCmd.AddCommand(commands.TodoCmd)

	// Add advanced command groups (for backwards compatibility and advanced features)
	rootCmd.AddCommand(commands.CommitCmd) // gitlike commit create/list/show
	rootCmd.AddCommand(commands.GitCmd)    // gitlike git (integration commands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
