package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"todo-cli/models"
)

// GitService handles integration with local Git repository
type GitService struct {
	repoPath string
}

// NewGitService creates a new Git service
func NewGitService() *GitService {
	// Try to find Git repository starting from current directory
	currentDir, _ := os.Getwd()
	repoPath := findGitRepo(currentDir)

	return &GitService{
		repoPath: repoPath,
	}
}

// findGitRepo recursively searches for .git directory
func findGitRepo(dir string) string {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return dir
	}

	parent := filepath.Dir(dir)
	if parent == dir {
		return "" // Reached root directory
	}

	return findGitRepo(parent)
}

// IsGitRepo checks if we're in a Git repository
func (g *GitService) IsGitRepo() bool {
	return g.repoPath != ""
}

// GetRepoPath returns the Git repository path
func (g *GitService) GetRepoPath() string {
	return g.repoPath
}

// GetCurrentBranch returns the current Git branch
func (g *GitService) GetCurrentBranch() (string, error) {
	if !g.IsGitRepo() {
		return "", fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// GetAllBranches returns all Git branches
func (g *GitService) GetAllBranches() ([]string, error) {
	if !g.IsGitRepo() {
		return nil, fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "branch", "-a")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var branches []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Remove markers and prefixes
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimPrefix(line, "  ")
		line = strings.TrimPrefix(line, "remotes/origin/")

		// Skip HEAD pointer
		if strings.Contains(line, "HEAD ->") {
			continue
		}

		// Avoid duplicates
		found := false
		for _, existing := range branches {
			if existing == line {
				found = true
				break
			}
		}
		if !found {
			branches = append(branches, line)
		}
	}

	return branches, nil
}

// GetRecentCommits returns recent Git commits
func (g *GitService) GetRecentCommits(limit int) ([]models.Commit, error) {
	if !g.IsGitRepo() {
		return nil, fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", limit), "--pretty=format:%H|%s|%an|%at")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []models.Commit
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			continue
		}

		// Get current branch for this commit
		branch, _ := g.GetCurrentBranch()

		commit := models.Commit{
			ID:        parts[0][:8], // Short hash
			Message:   parts[1],
			Author:    parts[2],
			Branch:    branch,
			CreatedAt: time.Unix(timestamp, 0),
			Todos:     []int{}, // Will be populated by analyzing commit
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// CheckoutBranch switches to a Git branch
func (g *GitService) CheckoutBranch(branchName string) error {
	if !g.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = g.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout failed: %s", string(output))
	}

	return nil
}

// CreateBranch creates a new Git branch
func (g *GitService) CreateBranch(branchName string) error {
	if !g.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = g.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git branch creation failed: %s", string(output))
	}

	return nil
}

// GetChangedFiles returns files changed in the working directory
func (g *GitService) GetChangedFiles() ([]string, error) {
	if !g.IsGitRepo() {
		return nil, fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 3 {
			// Format: "XY filename" where XY are status codes
			filename := line[3:]
			files = append(files, filename)
		}
	}

	return files, nil
}

// GetRemoteURL returns the Git remote URL
func (g *GitService) GetRemoteURL() (string, error) {
	if !g.IsGitRepo() {
		return "", fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// CommitChanges creates a Git commit
func (g *GitService) CommitChanges(message string) error {
	if !g.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	// Add all changes
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = g.repoPath
	_, err := addCmd.Output()
	if err != nil {
		return fmt.Errorf("git add failed: %v", err)
	}

	// Commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = g.repoPath
	output, err := commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %s", string(output))
	}

	return nil
}

// PushToRemote pushes changes to remote repository
func (g *GitService) PushToRemote() error {
	if !g.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "push")
	cmd.Dir = g.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push failed: %s", string(output))
	}

	return nil
}

// PushToRemoteWithDetails pushes changes and returns detailed information
func (g *GitService) PushToRemoteWithDetails() (string, error) {
	if !g.IsGitRepo() {
		return "", fmt.Errorf("not in a git repository")
	}

	// Get current branch
	currentBranch, err := g.GetCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}

	// Check if remote tracking branch exists
	trackingCmd := exec.Command("git", "rev-parse", "--abbrev-ref", currentBranch+"@{upstream}")
	trackingCmd.Dir = g.repoPath
	_, trackingErr := trackingCmd.Output()

	var pushCmd *exec.Cmd
	if trackingErr != nil {
		// No upstream branch, set it up
		pushCmd = exec.Command("git", "push", "--set-upstream", "origin", currentBranch)
	} else {
		// Upstream exists, normal push
		pushCmd = exec.Command("git", "push")
	}

	pushCmd.Dir = g.repoPath
	output, err := pushCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git push failed: %s", string(output))
	}

	return string(output), nil
}

// GetCommitsSinceLastPush returns commits that haven't been pushed
func (g *GitService) GetCommitsSinceLastPush() ([]models.Commit, error) {
	if !g.IsGitRepo() {
		return nil, fmt.Errorf("not in a git repository")
	}

	currentBranch, err := g.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	// Try to get commits since last push
	cmd := exec.Command("git", "log", "origin/"+currentBranch+"..HEAD", "--pretty=format:%H|%s|%an|%at")
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	if err != nil {
		// If origin/branch doesn't exist, get all commits
		cmd = exec.Command("git", "log", "HEAD", "--pretty=format:%H|%s|%an|%at")
		cmd.Dir = g.repoPath
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
	}

	var commits []models.Commit
	if len(output) == 0 {
		return commits, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			continue
		}

		commit := models.Commit{
			ID:        parts[0][:8],
			Message:   parts[1],
			Author:    parts[2],
			Branch:    currentBranch,
			CreatedAt: time.Unix(timestamp, 0),
			Todos:     []int{},
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// HasUnpushedChanges checks if there are unpushed commits
func (g *GitService) HasUnpushedChanges() (bool, error) {
	if !g.IsGitRepo() {
		return false, fmt.Errorf("not in a git repository")
	}

	commits, err := g.GetCommitsSinceLastPush()
	if err != nil {
		return false, err
	}

	return len(commits) > 0, nil
}

// GetPushStatus returns detailed push status information
func (g *GitService) GetPushStatus() (map[string]interface{}, error) {
	if !g.IsGitRepo() {
		return nil, fmt.Errorf("not in a git repository")
	}

	status := make(map[string]interface{})

	// Current branch
	currentBranch, err := g.GetCurrentBranch()
	if err != nil {
		return nil, err
	}
	status["current_branch"] = currentBranch

	// Check if there are unpushed commits
	unpushedCommits, err := g.GetCommitsSinceLastPush()
	if err == nil {
		status["unpushed_commits"] = len(unpushedCommits)
		status["commits"] = unpushedCommits
	}

	// Check if there are uncommitted changes
	changedFiles, err := g.GetChangedFiles()
	if err == nil {
		status["uncommitted_changes"] = len(changedFiles)
		status["changed_files"] = changedFiles
	}

	// Check remote URL
	remoteURL, err := g.GetRemoteURL()
	if err == nil {
		status["remote_url"] = remoteURL
	}

	return status, nil
}

// PullFromRemote pulls changes from remote repository
func (g *GitService) PullFromRemote() error {
	if !g.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = g.repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %s", string(output))
	}

	return nil
}

// SyncWithGit synchronizes todo repository with Git state
func (g *GitService) SyncWithGit(repo *models.Repository) error {
	if !g.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	// Get current Git branch
	currentBranch, err := g.GetCurrentBranch()
	if err != nil {
		return err
	}

	// Update current branch in todo repo
	repo.CurrentBranch = currentBranch

	// Get all Git branches and ensure they exist in todo repo
	gitBranches, err := g.GetAllBranches()
	if err != nil {
		return err
	}

	for _, gitBranch := range gitBranches {
		// Check if branch exists in todo repo
		found := false
		for _, todoBranch := range repo.Branches {
			if todoBranch.Name == gitBranch {
				found = true
				break
			}
		}

		// Create todo branch if it doesn't exist
		if !found {
			newBranch := models.Branch{
				Name:      gitBranch,
				CreatedAt: time.Now(),
				IsActive:  gitBranch == currentBranch,
				Todos:     []models.Todo{},
			}
			repo.Branches = append(repo.Branches, newBranch)
		}
	}

	return nil
}
