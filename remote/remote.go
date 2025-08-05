package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlike/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// RemoteService handles remote repository operations
type RemoteService struct {
	client *http.Client
}

// NewRemoteService creates a new remote service
func NewRemoteService() *RemoteService {
	return &RemoteService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PushRepository pushes the repository to a remote server
func (r *RemoteService) PushRepository(remote models.Remote, repo *models.Repository) error {
	switch remote.Type {
	case "http":
		return r.pushHTTP(remote, repo)
	case "file":
		return r.pushFile(remote, repo)
	default:
		return fmt.Errorf("unsupported remote type: %s", remote.Type)
	}
}

// PullRepository pulls the repository from a remote server
func (r *RemoteService) PullRepository(remote models.Remote) (*models.Repository, error) {
	switch remote.Type {
	case "http":
		return r.pullHTTP(remote)
	case "file":
		return r.pullFile(remote)
	default:
		return nil, fmt.Errorf("unsupported remote type: %s", remote.Type)
	}
}

// pushHTTP pushes repository to HTTP server
func (r *RemoteService) pushHTTP(remote models.Remote, repo *models.Repository) error {
	data, err := json.Marshal(repo)
	if err != nil {
		return fmt.Errorf("failed to marshal repository: %w", err)
	}

	req, err := http.NewRequest("POST", remote.URL+"/push", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication if credentials are available
	if username := os.Getenv("TODO_CLI_USERNAME"); username != "" {
		if password := os.Getenv("TODO_CLI_PASSWORD"); password != "" {
			req.SetBasicAuth(username, password)
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to push to remote: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error: %s", string(body))
	}

	return nil
}

// pullHTTP pulls repository from HTTP server
func (r *RemoteService) pullHTTP(remote models.Remote) (*models.Repository, error) {
	req, err := http.NewRequest("GET", remote.URL+"/pull", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if credentials are available
	if username := os.Getenv("TODO_CLI_USERNAME"); username != "" {
		if password := os.Getenv("TODO_CLI_PASSWORD"); password != "" {
			req.SetBasicAuth(username, password)
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to pull from remote: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	var repo models.Repository
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode repository: %w", err)
	}

	return &repo, nil
}

// pushFile pushes repository to file system
func (r *RemoteService) pushFile(remote models.Remote, repo *models.Repository) error {
	// Ensure directory exists
	dir := filepath.Dir(remote.URL)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(repo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal repository: %w", err)
	}

	err = os.WriteFile(remote.URL, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// pullFile pulls repository from file system
func (r *RemoteService) pullFile(remote models.Remote) (*models.Repository, error) {
	data, err := os.ReadFile(remote.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var repo models.Repository
	err = json.Unmarshal(data, &repo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal repository: %w", err)
	}

	return &repo, nil
}

// MergeRepositories merges a remote repository with local repository
func (r *RemoteService) MergeRepositories(local, remote *models.Repository) *models.Repository {
	merged := *local // Start with local copy

	// Merge branches
	for _, remoteBranch := range remote.Branches {
		found := false
		for i, localBranch := range merged.Branches {
			if localBranch.Name == remoteBranch.Name {
				// Merge todos from remote branch
				merged.Branches[i] = r.mergeBranches(localBranch, remoteBranch)
				found = true
				break
			}
		}
		if !found {
			// Add new branch from remote
			merged.Branches = append(merged.Branches, remoteBranch)
		}
	}

	// Merge commits (avoid duplicates)
	commitMap := make(map[string]bool)
	for _, commit := range merged.Commits {
		commitMap[commit.ID] = true
	}

	for _, remoteCommit := range remote.Commits {
		if !commitMap[remoteCommit.ID] {
			merged.Commits = append(merged.Commits, remoteCommit)
		}
	}

	// Update next todo ID to avoid conflicts
	if remote.NextTodoID > merged.NextTodoID {
		merged.NextTodoID = remote.NextTodoID
	}

	merged.LastSync = time.Now()

	return &merged
}

// mergeBranches merges todos from two branches
func (r *RemoteService) mergeBranches(local, remote models.Branch) models.Branch {
	merged := local
	todoMap := make(map[int]bool)

	// Track existing todos
	for _, todo := range merged.Todos {
		todoMap[todo.ID] = true
	}

	// Add new todos from remote
	for _, remoteTodo := range remote.Todos {
		if !todoMap[remoteTodo.ID] {
			merged.Todos = append(merged.Todos, remoteTodo)
		} else {
			// Update existing todo if remote is newer
			for i, localTodo := range merged.Todos {
				if localTodo.ID == remoteTodo.ID && remoteTodo.UpdatedAt.After(localTodo.UpdatedAt) {
					merged.Todos[i] = remoteTodo
					break
				}
			}
		}
	}

	return merged
}
