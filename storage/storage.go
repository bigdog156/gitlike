package storage

import (
	"encoding/json"
	"fmt"
	"gitlike/models"
	"os"
	"path/filepath"
	"time"
)

const (
	dataDir  = ".tododata"
	repoFile = "repository.json"
)

// Storage handles data persistence
type Storage struct {
	dataPath string
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	homeDir, _ := os.UserHomeDir()
	dataPath := filepath.Join(homeDir, dataDir)

	// Create data directory if it doesn't exist
	os.MkdirAll(dataPath, 0755)

	return &Storage{
		dataPath: dataPath,
	}
}

// LoadRepository loads the repository from disk
func (s *Storage) LoadRepository() (*models.Repository, error) {
	repoPath := filepath.Join(s.dataPath, repoFile)

	// If file doesn't exist, create a new repository
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
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
				Enabled:        false,
				AutoSync:       false,
				RepoPath:       "",
				RemoteURL:      "",
				LastGitSync:    time.Time{},
				AutoCommit:     false,
				CommitTemplate: "todo: {{.Message}}",
			},
		}
		s.SaveRepository(repo)
		return repo, nil
	}

	data, err := os.ReadFile(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository: %w", err)
	}

	var repo models.Repository
	err = json.Unmarshal(data, &repo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repository: %w", err)
	}

	return &repo, nil
}

// SaveRepository saves the repository to disk
func (s *Storage) SaveRepository(repo *models.Repository) error {
	repoPath := filepath.Join(s.dataPath, repoFile)

	data, err := json.MarshalIndent(repo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal repository: %w", err)
	}

	err = os.WriteFile(repoPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write repository: %w", err)
	}

	return nil
}

// GetCurrentBranch returns the current active branch
func (s *Storage) GetCurrentBranch(repo *models.Repository) *models.Branch {
	for i := range repo.Branches {
		if repo.Branches[i].Name == repo.CurrentBranch {
			return &repo.Branches[i]
		}
	}
	return nil
}

// GetBranchByName returns a branch by name
func (s *Storage) GetBranchByName(repo *models.Repository, name string) *models.Branch {
	for i := range repo.Branches {
		if repo.Branches[i].Name == name {
			return &repo.Branches[i]
		}
	}
	return nil
}
