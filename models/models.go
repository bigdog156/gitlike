package models

import (
	"time"
)

// Todo represents a task
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // "pending", "in-progress", "completed"
	Priority    string    `json:"priority"` // "low", "medium", "high"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	BranchName  string    `json:"branch_name"`
}

// Branch represents a development branch
type Branch struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
	Todos     []Todo    `json:"todos"`
}

// Commit represents a commit with todos
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Branch    string    `json:"branch"`
	Todos     []int     `json:"todos"` // Todo IDs included in this commit
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
}

// Remote represents a remote repository
type Remote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"` // "http", "file", "git"
}

// Repository represents the entire todo repository
type Repository struct {
	Branches      []Branch  `json:"branches"`
	Commits       []Commit  `json:"commits"`
	CurrentBranch string    `json:"current_branch"`
	NextTodoID    int       `json:"next_todo_id"`
	Remotes       []Remote  `json:"remotes"`
	LastSync      time.Time `json:"last_sync"`
}
