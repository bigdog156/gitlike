package commands

import (
	"fmt"
	"gitlike/models"
	"gitlike/remote"

	"github.com/spf13/cobra"
)

var remoteService = remote.NewRemoteService()

var RemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote repository commands",
}

var remoteAddCmd = &cobra.Command{
	Use:   "add [name] [url]",
	Short: "Add a remote repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		url := args[1]

		remoteType, _ := cmd.Flags().GetString("type")

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Check if remote already exists
		for _, remote := range repo.Remotes {
			if remote.Name == name {
				fmt.Printf("Remote '%s' already exists\n", name)
				return
			}
		}

		// Add new remote
		newRemote := models.Remote{
			Name: name,
			URL:  url,
			Type: remoteType,
		}

		repo.Remotes = append(repo.Remotes, newRemote)

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Added remote '%s': %s (%s)\n", name, url, remoteType)
	},
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List remote repositories",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		if len(repo.Remotes) == 0 {
			fmt.Println("No remotes configured")
			return
		}

		fmt.Println("Remotes:")
		for _, remote := range repo.Remotes {
			fmt.Printf("  %s\t%s (%s)\n", remote.Name, remote.URL, remote.Type)
		}
	},
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a remote repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Find and remove remote
		found := false
		for i, remote := range repo.Remotes {
			if remote.Name == name {
				repo.Remotes = append(repo.Remotes[:i], repo.Remotes[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Remote '%s' not found\n", name)
			return
		}

		err = storage_instance.SaveRepository(repo)
		if err != nil {
			fmt.Printf("Error saving repository: %v\n", err)
			return
		}

		fmt.Printf("Removed remote '%s'\n", name)
	},
}

var PushCmd = &cobra.Command{
	Use:   "push [remote]",
	Short: "Push commits to remote repository",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := "origin"
		if len(args) > 0 {
			remoteName = args[0]
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Find remote
		var targetRemote *models.Remote
		for _, remote := range repo.Remotes {
			if remote.Name == remoteName {
				targetRemote = &remote
				break
			}
		}

		if targetRemote == nil {
			fmt.Printf("Remote '%s' not found\n", remoteName)
			return
		}

		fmt.Printf("Pushing to %s (%s)...\n", targetRemote.Name, targetRemote.URL)

		err = remoteService.PushRepository(*targetRemote, repo)
		if err != nil {
			fmt.Printf("Push failed: %v\n", err)
			return
		}

		fmt.Printf("Successfully pushed to %s\n", targetRemote.Name)
	},
}

var PullCmd = &cobra.Command{
	Use:   "pull [remote]",
	Short: "Pull and merge changes from remote repository",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := "origin"
		if len(args) > 0 {
			remoteName = args[0]
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Find remote
		var targetRemote *models.Remote
		for _, remote := range repo.Remotes {
			if remote.Name == remoteName {
				targetRemote = &remote
				break
			}
		}

		if targetRemote == nil {
			fmt.Printf("Remote '%s' not found\n", remoteName)
			return
		}

		fmt.Printf("Pulling from %s (%s)...\n", targetRemote.Name, targetRemote.URL)

		remoteRepo, err := remoteService.PullRepository(*targetRemote)
		if err != nil {
			fmt.Printf("Pull failed: %v\n", err)
			return
		}

		// Merge remote changes
		mergedRepo := remoteService.MergeRepositories(repo, remoteRepo)

		err = storage_instance.SaveRepository(mergedRepo)
		if err != nil {
			fmt.Printf("Error saving merged repository: %v\n", err)
			return
		}

		fmt.Printf("Successfully pulled and merged from %s\n", targetRemote.Name)
		fmt.Printf("- %d branches synced\n", len(remoteRepo.Branches))
		fmt.Printf("- %d commits synced\n", len(remoteRepo.Commits))
	},
}

var FetchCmd = &cobra.Command{
	Use:   "fetch [remote]",
	Short: "Fetch changes from remote repository without merging",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := "origin"
		if len(args) > 0 {
			remoteName = args[0]
		}

		repo, err := storage_instance.LoadRepository()
		if err != nil {
			fmt.Printf("Error loading repository: %v\n", err)
			return
		}

		// Find remote
		var targetRemote *models.Remote
		for _, remote := range repo.Remotes {
			if remote.Name == remoteName {
				targetRemote = &remote
				break
			}
		}

		if targetRemote == nil {
			fmt.Printf("Remote '%s' not found\n", remoteName)
			return
		}

		fmt.Printf("Fetching from %s (%s)...\n", targetRemote.Name, targetRemote.URL)

		remoteRepo, err := remoteService.PullRepository(*targetRemote)
		if err != nil {
			fmt.Printf("Fetch failed: %v\n", err)
			return
		}

		// Show what would be merged
		fmt.Printf("Remote has:\n")
		fmt.Printf("- %d branches\n", len(remoteRepo.Branches))
		fmt.Printf("- %d commits\n", len(remoteRepo.Commits))
		fmt.Printf("- Last sync: %s\n", remoteRepo.LastSync.Format("2006-01-02 15:04:05"))
		fmt.Println("\nUse 'todo pull' to merge these changes")
	},
}

var SyncCmd = &cobra.Command{
	Use:   "sync [remote]",
	Short: "Synchronize with remote (pull then push)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := "origin"
		if len(args) > 0 {
			remoteName = args[0]
		}

		fmt.Printf("Synchronizing with %s...\n", remoteName)

		// First pull
		PullCmd.Run(cmd, args)

		// Then push
		PushCmd.Run(cmd, args)

		fmt.Printf("Synchronization with %s complete\n", remoteName)
	},
}

func init() {
	// Add flags
	remoteAddCmd.Flags().StringP("type", "t", "http", "Remote type (http, file)")

	// Add subcommands
	RemoteCmd.AddCommand(remoteAddCmd)
	RemoteCmd.AddCommand(remoteListCmd)
	RemoteCmd.AddCommand(remoteRemoveCmd)

	// Add standalone commands that will be added to root
}
