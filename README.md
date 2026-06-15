# Task Manager (Go + HTML/CSS/JS)

A simple task manager app with a Go (`net/http`) backend and a vanilla JavaScript frontend.

## Features
- Add tasks
- List tasks
- Toggle task completion (done/pending)
- View a task by ID
- Delete tasks

## Go Concepts Used

- net/http for building the web server
- http.NewServeMux for routing requests
- HTTP handlers using http.HandleFunc
- REST API methods (GET, POST, PUT, DELETE)
- JSON handling using encoding/json
- Structs for data modeling (Task struct)
- In-memory data storage using slices ([]Task)
- Middleware for CORS handling
- Goroutines for background tasks
- Channels for communication between goroutines
- sync.Mutex for safe concurrent access
- URL query parsing using r.URL.Query()
- Logging using log package
  
## Project Structure
- `main.go` - Go web server + REST API endpoints + static file hosting
- `static/index.html` - Task Manager UI
- `static/app.js` - Frontend logic (fetch calls + DOM updates)
- `static/styles.css` - Styling

## How it works
### Backend (Go)
The server:
- Serves static files from the `static/` directory
- Exposes JSON REST endpoints under `/api/`

Implemented endpoints:
- `GET    /api/tasks` - returns `{ "tasks": [...] }`
- `POST   /api/add` - accepts `{ "title": "..." }` and returns `{ "task": {...} }`
- `GET    /api/task?id=1` - returns `{ "task": {...} }` or `{ "error": "task not found" }`
- `PUT    /api/done?id=1` - toggles `done` and returns the updated task
- `DELETE /api/delete?id=1` - deletes a task and returns the deleted task

Data is stored **in-memory** (restart clears tasks).

### Frontend (Browser)
`static/app.js`:
- Fetches `/api/tasks` on load
- Renders each task with buttons:
  - Toggle done via `PUT /api/done?id=...`
  - Delete via `DELETE /api/delete?id=...`
- Add task form sends `POST /api/add` with JSON payload `{ title }`

## Tech Stack
- **Backend:** Go (`net/http`), in-memory state, JSON
- **Frontend:** HTML + CSS + vanilla JavaScript (Fetch API)

## Live Demo
🔗 https://go-project-production-023a.up.railway.app/

## Run the Project
From the project root directory `task-manager/`, run:

```bash
go run .

