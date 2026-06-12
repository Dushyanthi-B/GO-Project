package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

var (
	tasks     []Task
	nextTaskID = 1
	mu        sync.Mutex
)

func main() {
	// In-memory only; start empty.
	tasks = make([]Task, 0)

	staticDir := filepath.Join("static")
	// If you run from project root, staticDir exists.
	// If you run from elsewhere, try locating relative to executable.
	if _, err := os.Stat(staticDir); err != nil {
		exe, _ := os.Executable()
		staticDir = filepath.Join(filepath.Dir(exe), "static")
	}

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/tasks", getTasksHandler)
	apiMux.HandleFunc("/api/add", addTaskHandler)

	fileServer := http.FileServer(http.Dir(staticDir))
	// Serve static assets under "/".
	rootMux := http.NewServeMux()
	rootMux.Handle("/", fileServer)

	// Combine routers.
	mux := http.NewServeMux()
	mux.Handle("/api/", apiMux)
	mux.Handle("/", rootMux)

	addr := ":8001"
	log.Printf("Task Manager running at http://localhost%s", addr)
	log.Printf("Static files dir: %s", staticDir)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Same-origin normally makes this unnecessary, but it keeps it beginner-friendly.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}

type addTaskRequest struct {
	Text string `json:"text"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req addTaskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid json: %v", err)})
		return
	}

	text := req.Text
	if len(text) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text is required"})
		return
	}

	mu.Lock()
	task := Task{ID: nextTaskID, Text: text, Completed: false}
	nextTaskID++
	tasks = append(tasks, task)
	mu.Unlock()

	writeJSON(w, http.StatusCreated, map[string]any{"task": task})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Printf("failed to encode json: %v", err)
	}
}

