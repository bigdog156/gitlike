package commands

import (
	"fmt"
	"gitlike/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// GitLike Init Command - gitlike init
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty GitLike repository",
	Long:  "Initialize a new GitLike repository with todo tracking",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Git repository first
		gitInit := exec.Command("git", "init")
		err := gitInit.Run()
		if err != nil {
			fmt.Printf("Error initializing Git repository: %v\n", err)
			return
		}

		// Initialize GitLike repository
		repo := &models.Repository{
			Branches: []models.Branch{
				{
					Name:      "main",
					CreatedAt: time.Now(),
					IsActive:  true,
					Todos:     []models.Todo{},
				},
			},
			Commits:       []models.Commit{},
			CurrentBranch: "main",
			NextTodoID:    1,
			Remotes:       []models.Remote{},
			LastSync:      time.Time{},
			GitIntegration: models.GitConfig{
				Enabled:        true,
				AutoSync:       true,
				RepoPath:       "",
				RemoteURL:      "",
				LastGitSync:    time.Time{},
				AutoCommit:     false,
				CommitTemplate: "{{.Message}}",
			},
		}

		currentDir, _ := os.Getwd()
		repo.GitIntegration.RepoPath = currentDir

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Initialized empty GitLike repository in %s\n", filepath.Join(currentDir, ".git"))
		fmt.Println("GitLike todo tracking enabled")
	},
}

// GitLike Add Command - gitlike add
var AddCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add file contents to the index",
	Long:  "Add file contents to the staging area",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Nothing specified, nothing added.")
			return
		}

		// Execute git add with the provided arguments
		gitArgs := append([]string{"add"}, args...)
		gitAdd := exec.Command("git", gitArgs...)
		output, err := gitAdd.CombinedOutput()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			return
		}

		// Show what was added
		if len(args) == 1 && args[0] == "." {
			fmt.Println("Added all changes to staging area")
		} else {
			fmt.Printf("Added %d file(s) to staging area:\n", len(args))
			for _, file := range args {
				fmt.Printf("  %s\n", file)
			}
		}
	},
}

// GitLike Status Command - gitlike status
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long:  "Show the working tree status with todo context",
	Run: func(cmd *cobra.Command, args []string) {
		// Get Git status first
		gitStatus := exec.Command("git", "status", "--porcelain")
		output, err := gitStatus.Output()
		if err != nil {
			fmt.Printf("Error getting git status: %v\n", err)
			return
		}

		// Get current branch
		gitBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branchOutput, err := gitBranch.Output()
		currentGitBranch := "main"
		if err == nil {
			currentGitBranch = strings.TrimSpace(string(branchOutput))
		}

		fmt.Printf("On branch %s\n", currentGitBranch)

		// Load GitLike repository for todo context
		repo, err := storage_instance.LoadRepository()
		if err == nil {
			currentBranch := storage_instance.GetCurrentBranch(repo)
			if currentBranch != nil {
				activeTodos := 0
				completedTodos := 0
				for _, todo := range currentBranch.Todos {
					if todo.IsActive {
						activeTodos++
					}
					if todo.Status == "completed" {
						completedTodos++
					}
				}
				
				if activeTodos > 0 || completedTodos > 0 {
					fmt.Printf("Todo status: %d active, %d completed\n", activeTodos, completedTodos)
				}
			}
		}

		// Parse and display Git status
		statusLines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(statusLines) == 1 && statusLines[0] == "" {
			fmt.Println("nothing to commit, working tree clean")
			return
		}

		var staged []string
		var modified []string
		var untracked []string

		for _, line := range statusLines {
			if len(line) < 3 {
				continue
			}
			
			statusCode := line[:2]
			fileName := line[3:]

			switch statusCode {
			case "A ", "M ", "D ", "R ", "C ":
				staged = append(staged, fileName)
			case " M", " D":
				modified = append(modified, fileName)
			case "??":
				untracked = append(untracked, fileName)
			case "AM", "MM":
				staged = append(staged, fileName)
				modified = append(modified, fileName)
			}
		}

		if len(staged) > 0 {
			fmt.Println("\nChanges to be committed:")
			fmt.Println("  (use \"gitlike reset HEAD <file>...\" to unstage)")
			fmt.Println()
			for _, file := range staged {
				fmt.Printf("\tmodified:   %s\n", file)
			}
		}

		if len(modified) > 0 {
			fmt.Println("\nChanges not staged for commit:")
			fmt.Println("  (use \"gitlike add <file>...\" to update what will be committed)")
			fmt.Println("  (use \"gitlike checkout -- <file>...\" to discard changes in working directory)")
			fmt.Println()
			for _, file := range modified {
				fmt.Printf("\tmodified:   %s\n", file)
			}
		}

		if len(untracked) > 0 {
			fmt.Println("\nUntracked files:")
			fmt.Println("  (use \"gitlike add <file>...\" to include in what will be committed)")
			fmt.Println()
			for _, file := range untracked {
				fmt.Printf("\t%s\n", file)
			}
		}
	},
}

// GitLike Commit Command - gitlike commit
var CommitDirectCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository with todo tracking",
	Long:  "Create a new commit with changes staged in the index, automatically tracking related todos",
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		if message == "" && len(args) > 0 {
			message = strings.Join(args, " ")
		}
		
		if message == "" {
			fmt.Println("Error: commit message required")
			fmt.Println("Usage: gitlike commit -m \"your message\"")
			return
		}

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

		// Find completed todos and active todo for GitLike tracking
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

		// Determine todos to include in commit
		var todosToInclude []int
		todoContext := ""
		
		if activeTodo != nil {
			todosToInclude = append(todosToInclude, activeTodo.ID)
			todoContext = fmt.Sprintf(" (working on #%d: %s)", activeTodo.ID, activeTodo.Title)
		}
		
		if len(completedTodos) > 0 {
			todosToInclude = append(todosToInclude, completedTodos...)
			if activeTodo != nil {
				todoContext += fmt.Sprintf(" + %d completed", len(completedTodos))
			} else {
				todoContext = fmt.Sprintf(" (completed #%v)", completedTodos)
			}
		}

		// Create Git commit first
		gitCommit := exec.Command("git", "commit", "-m", message)
		output, err := gitCommit.CombinedOutput()
		if err != nil {
			fmt.Printf("Git commit failed: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			return
		}

		// Parse git commit output to get commit hash
		outputStr := string(output)
		gitCommitHash := ""
		if strings.Contains(outputStr, "[") && strings.Contains(outputStr, "]") {
			// Extract commit hash from output like "[main abc1234] commit message"
			parts := strings.Split(outputStr, "]")
			if len(parts) > 0 {
				leftPart := parts[0]
				if idx := strings.LastIndex(leftPart, " "); idx != -1 {
					gitCommitHash = leftPart[idx+1:]
				}
			}
		}

		// Create GitLike commit record
		if len(todosToInclude) > 0 {
			// Get current user
			author := "unknown"
			if gitConfig := exec.Command("git", "config", "user.name"); gitConfig.Run() == nil {
				if authorOutput, err := gitConfig.Output(); err == nil {
					author = strings.TrimSpace(string(authorOutput))
				}
			}

			commit := models.Commit{
				ID:        gitCommitHash,
				Message:   message,
				Branch:    currentBranch.Name,
				Todos:     todosToInclude,
				CreatedAt: time.Now(),
				Author:    author,
			}

			if activeTodo != nil {
				commit.ActiveTodo = &activeTodo.ID
			}

			// Link commit to todos
			for i := range repo.Branches {
				if repo.Branches[i].Name == currentBranch.Name {
					for j := range repo.Branches[i].Todos {
						for _, todoID := range todosToInclude {
							if repo.Branches[i].Todos[j].ID == todoID {
								if repo.Branches[i].Todos[j].Commits == nil {
									repo.Branches[i].Todos[j].Commits = []string{}
								}
								repo.Branches[i].Todos[j].Commits = append(repo.Branches[i].Todos[j].Commits, commit.ID)
								break
							}
						}
					}
					break
				}
			}

			repo.Commits = append(repo.Commits, commit)
			storage_instance.SaveRepository(repo)

			fmt.Printf("✅ Git commit successful%s\n", todoContext)
		} else {
			fmt.Printf("✅ Git commit successful\n")
		}

		// Show git output
		fmt.Printf("%s", string(output))
	},
}

// GitLike Log Command - gitlike log
var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs with todo history",
	Long:  "Display commit history with associated todo information",
	Run: func(cmd *cobra.Command, args []string) {
		// Get GitLike commits first
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Display GitLike commits with todo context
		fmt.Println("GitLike commit history with todo tracking:\n")
		
		for i := len(repo.Commits) - 1; i >= 0; i-- {
			commit := repo.Commits[i]
			fmt.Printf("commit %s\n", commit.ID)
			fmt.Printf("Author: %s\n", commit.Author)
			fmt.Printf("Date: %s\n", commit.CreatedAt.Format("Mon Jan 2 15:04:05 2006 -0700"))
			fmt.Printf("\n    %s\n", commit.Message)
			
			if len(commit.Todos) > 0 {
				fmt.Printf("    \n    Todos included: ")
				for j, todoID := range commit.Todos {
					if j > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("#%d", todoID)
				}
				fmt.Printf("\n")
			}
			
			if commit.ActiveTodo != nil {
				fmt.Printf("    Active todo: #%d\n", *commit.ActiveTodo)
			}
			fmt.Println()
		}

		// Also show git log for comparison
		oneline, _ := cmd.Flags().GetBool("oneline")
		if oneline {
			gitLog := exec.Command("git", "log", "--oneline")
			output, err := gitLog.Output()
			if err == nil {
				fmt.Println("Git commit history:")
				fmt.Printf("%s", string(output))
			}
		}
	},
}

// GitLike Checkout Command - gitlike checkout
var CheckoutCmd = &cobra.Command{
	Use:   "checkout [branch|commit]",
	Short: "Switch branches or restore working tree files", 
	Long:  "Switch to another branch with todo synchronization",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: gitlike checkout <branch>")
			return
		}

		target := args[0]
		
		// Create new branch if -b flag is provided
		createBranch, _ := cmd.Flags().GetBool("b")
		if createBranch {
			// Create Git branch
			gitCheckout := exec.Command("git", "checkout", "-b", target)
			output, err := gitCheckout.CombinedOutput()
			if err != nil {
				fmt.Printf("Error creating branch: %v\n", err)
				fmt.Printf("Output: %s\n", string(output))
				return
			}

			// Create GitLike branch
			repo, err := storage_instance.LoadRepository()
			if err != nil {
				fmt.Printf("Error loading repository: %v\n", err)
				return
			}

			// Add new branch to GitLike
			newBranch := models.Branch{
				Name:      target,
				CreatedAt: time.Now(),
				IsActive:  true,
				Todos:     []models.Todo{},
			}

			// Deactivate current branch
			for i := range repo.Branches {
				repo.Branches[i].IsActive = false
			}

			repo.Branches = append(repo.Branches, newBranch)
			repo.CurrentBranch = target

			err = storage_instance.SaveRepository(repo)
			if err != nil {
				fmt.Printf("Error saving repository: %v\n", err)
				return
			}

			fmt.Printf("Switched to a new branch '%s'\n", target)
			return
		}

		// Switch to existing branch
		gitCheckout := exec.Command("git", "checkout", target)
		output, err := gitCheckout.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			return
		}

		// Update GitLike branch
		repo, err := storage_instance.LoadRepository()
		if err == nil {
			// Deactivate all branches
			for i := range repo.Branches {
				repo.Branches[i].IsActive = false
			}

			// Activate target branch or create if not exists
			found := false
			for i := range repo.Branches {
				if repo.Branches[i].Name == target {
					repo.Branches[i].IsActive = true
					found = true
					break
				}
			}

			if !found {
				// Create new GitLike branch for existing Git branch
				newBranch := models.Branch{
					Name:      target,
					CreatedAt: time.Now(),
					IsActive:  true,
					Todos:     []models.Todo{},
				}
				repo.Branches = append(repo.Branches, newBranch)
			}

			repo.CurrentBranch = target
			storage_instance.SaveRepository(repo)
		}

		fmt.Printf("Switched to branch '%s'\n", target)
	},
}

func init() {
	// Add flags for various commands
	CommitDirectCmd.Flags().StringP("message", "m", "", "Commit message")
	LogCmd.Flags().BoolP("oneline", "", false, "Show one line per commit")
	CheckoutCmd.Flags().BoolP("b", "b", false, "Create and switch to new branch")
}
