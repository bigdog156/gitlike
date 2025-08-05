package main

import (
	"encoding/json"
	"fmt"
	"gitlike/models"
	"io"
	"log"
	"net/http"
	"os"
)

const dataFile = "server_repository.json"

type Server struct {
	dataPath string
}

func NewServer() *Server {
	return &Server{
		dataPath: dataFile,
	}
}

func (s *Server) loadRepository() (*models.Repository, error) {
	if _, err := os.Stat(s.dataPath); os.IsNotExist(err) {
		// Return empty repository if file doesn't exist
		return &models.Repository{
			Branches:      []models.Branch{},
			Commits:       []models.Commit{},
			CurrentBranch: "main",
			NextTodoID:    1,
			Remotes:       []models.Remote{},
		}, nil
	}

	data, err := os.ReadFile(s.dataPath)
	if err != nil {
		return nil, err
	}

	var repo models.Repository
	err = json.Unmarshal(data, &repo)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (s *Server) saveRepository(repo *models.Repository) error {
	data, err := json.MarshalIndent(repo, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.dataPath, data, 0644)
}

func (s *Server) handlePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var clientRepo models.Repository
	err = json.Unmarshal(body, &clientRepo)
	if err != nil {
		http.Error(w, "Failed to parse repository", http.StatusBadRequest)
		return
	}

	// Save the pushed repository
	err = s.saveRepository(&clientRepo)
	if err != nil {
		http.Error(w, "Failed to save repository", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received push: %d branches, %d commits\n", len(clientRepo.Branches), len(clientRepo.Commits))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Push successful"))
}

func (s *Server) handlePull(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repo, err := s.loadRepository()
	if err != nil {
		http.Error(w, "Failed to load repository", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)

	fmt.Printf("Served pull: %d branches, %d commits\n", len(repo.Branches), len(repo.Commits))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	repo, err := s.loadRepository()
	if err != nil {
		http.Error(w, "Failed to load repository", http.StatusInternalServerError)
		return
	}

	status := map[string]interface{}{
		"branches":       len(repo.Branches),
		"commits":        len(repo.Commits),
		"current_branch": repo.CurrentBranch,
		"last_sync":      repo.LastSync,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewServer()

	http.HandleFunc("/push", server.handlePush)
	http.HandleFunc("/pull", server.handlePull)
	http.HandleFunc("/status", server.handleStatus)

	// Serve static files for web interface (optional)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Todo CLI Server</title>
</head>
<body>
    <h1>Todo CLI Remote Server</h1>
    <p>Endpoints:</p>
    <ul>
        <li><a href="/status">GET /status</a> - Repository status</li>
        <li>POST /push - Push repository</li>
        <li>GET /pull - Pull repository</li>
    </ul>
</body>
</html>
			`)
		} else {
			http.NotFound(w, r)
		}
	})

	fmt.Printf("Todo CLI server starting on port %s\n", port)
	fmt.Printf("Access the server at: http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
