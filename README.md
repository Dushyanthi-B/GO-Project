# Task Manager (Go + HTML/CSS/JS)

A simple task manager app with a Go (`net/http`) backend and a vanilla JavaScript frontend.

## Features
- Add tasks
- List tasks
- Toggle task completion (done/pending)
- View a task by ID
- Delete tasks

## Project Structure
- `main.go` - Go web server + REST API endpoints + static file hosting
- `static/index.html` - Task Manager UI
- `static/app.js` - Frontend logic (fetch calls + DOM updates)
- `static/styles.css` - Styling
- `static/home.html` - Home page that redirects to the Task Manager UI

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

## Run
From `task-manager/`:
```bash
go run .
```
Then open:
- Task Manager UI: `http://localhost:8001/index.html`
- Home page (redirects): `http://localhost:8001/`

## Notes
- The server listens on port **8002**.
- Since storage is in-memory, tasks will be lost when the server stops.

