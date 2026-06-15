package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks      []Task
	nextTaskID = 1
	mu         sync.Mutex
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
	apiMux.HandleFunc("/api/task", getTaskHandler)
	apiMux.HandleFunc("/api/done", toggleDoneHandler)
	apiMux.HandleFunc("/api/delete", deleteTaskHandler)

	fileServer := http.FileServer(http.Dir(staticDir))
	// Serve static assets under "/".
	rootMux := http.NewServeMux()
	rootMux.Handle("/", fileServer)

	// Combine routers.
	mux := http.NewServeMux()
	mux.Handle("/api/", apiMux)
	mux.Handle("/", rootMux)

	// Background worker: goroutines + channel.
	// One goroutine reads log messages from logCh and prints them.

	type logEvent struct{ msg string }
	logCh := make(chan logEvent, 32)
	go func() {
		for ev := range logCh {
			log.Println(ev.msg)
		}
	}()

	go func() {
		for {
			logCh <- logEvent{msg: "Background worker alive..."}
			time.Sleep(10 * time.Second)
		}
	}()

	addr := ":8001"
	log.Printf("Task Manager running at http://localhost%s", addr)
	log.Printf("Static files dir: %s", staticDir)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
	Title string `json:"title"`
}

// This makes the API more tolerant to frontend/consumer variations.
type addTaskRequestCompat struct {
	Title string `json:"title"`
	Task  string `json:"task"`
	Name  string `json:"name"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req addTaskRequest
	dec := json.NewDecoder(r.Body)
	// Be tolerant to unknown fields to avoid breaking consumers.
	if err := dec.Decode(&req); err != nil {
		// Backward/forward compat fallback
		var compat addTaskRequestCompat
		_ = json.NewDecoder(r.Body).Decode(&compat)
		if compat.Title != "" {
			title := compat.Title
			// continue below
			mu.Lock()
			task := Task{ID: nextTaskID, Title: title, Done: false}
			nextTaskID++
			tasks = append(tasks, task)
			mu.Unlock()
			writeJSON(w, http.StatusCreated, map[string]any{"task": task})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid json: %v", err)})
		return
	}

	title := req.Title

	if len(title) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	mu.Lock()
	task := Task{ID: nextTaskID, Title: title, Done: false}
	nextTaskID++
	tasks = append(tasks, task)
	mu.Unlock()

	writeJSON(w, http.StatusCreated, map[string]any{"task": task})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id, ok := parseIDQuery(w, r)
	if !ok {
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for _, t := range tasks {
		if t.ID == id {
			writeJSON(w, http.StatusOK, map[string]any{"task": t})
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
}

func toggleDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id, ok := parseIDQuery(w, r)
	if !ok {
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = !tasks[i].Done
			writeJSON(w, http.StatusOK, map[string]any{"task": tasks[i]})
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id, ok := parseIDQuery(w, r)
	if !ok {
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range tasks {
		if tasks[i].ID == id {
			deleted := tasks[i]
			tasks = append(tasks[:i], tasks[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]any{"task": deleted})
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
}

func parseIDQuery(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return 0, false
	}

	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return 0, false
	}
	return id, true
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
