# Task Manager CLI & HTTP Server in Go

A task management application built in **Go**, featuring both a **CLI mode** and an **HTTP server mode**. This project demonstrates core Go concepts including **structs, slices, error handling, multithreading, file I/O, JSON encoding/decoding, and HTTP servers**.

---

## Features

### CLI Mode
- Add a task: `add <task title>`  
- Remove a task: `remove <task id>`  
- List tasks: `list [--pending | --completed]`  
- Complete a task: `complete <task id>`  
- Remove all completed tasks: `update`  
- Exit the CLI: `exit` or `quit`

### Server Mode
- RESTful HTTP API endpoints:
  - `POST /tasks/add` → Add a task (JSON: `{ "title": "Buy milk" }`)  
  - `POST /tasks/remove` → Remove task (JSON: `{ "id": 1 }`)  
  - `POST /tasks/complete` → Complete task (JSON: `{ "id": 1 }`)  
  - `GET /tasks/list` → List tasks, optional query `?pending=true` or `?completed=true`  
  - `POST /tasks/update` → Remove all completed tasks, returns removed tasks  
  - `GET /logs` → Get the full log file  
  - `POST /exit` → Shut down server (JSON: `{ "password": "<EXIT_PASSWORD>" }`)  

- Auto-saves tasks every 10 seconds  
- Thread-safe logging and task updates  

---

## Concepts Demonstrated

### 1. **Multithreading / Concurrency**
- Auto-save and logging run in **separate goroutines**, so the main application flow is not blocked.  
- **Mutexes (`sync.Mutex`)** ensure safe access to shared resources like tasks and logs.  

### 2. **File I/O**
- Tasks and logs are persisted in JSON files.  
- Supports separate files for CLI mode and server mode for flexibility.  

### 3. **JSON Encoding/Decoding**
- Tasks and server requests/responses are encoded/decoded using `encoding/json`.  
- This allows structured data to be easily sent over HTTP or stored in files.
- Easy and flexible configuration

### 4. **HTTP Server**
- Implements RESTful endpoints with `net/http`.  
- Handlers use JSON requests and responses.  
- Graceful shutdown with `http.Server.Shutdown` allows stopping the server safely from a request (`/exit`).  

---

## Configuration

- Separate configuration files store settings such as task file names, log file names, and the exit password.  
- The program reads these config files on startup to initialize the server or CLI mode.  

---

## Running the Program

### CLI Mode
```bash
go run main.go

### SERVER Mode
go run main.go server

Example curl Requests (Server)

Add a task:
curl -X POST http://localhost:8080/tasks/add \
     -H "Content-Type: application/json" \
     -d '{"title":"Buy milk"}'


Complete a task:
curl -X POST http://localhost:8080/tasks/complete \
     -H "Content-Type: application/json" \
     -d '{"id":1}'


List all tasks:
curl http://localhost:8080/tasks/list


List only pending tasks:
curl http://localhost:8080/tasks/list?pending=true


Remove completed tasks:
curl -X POST http://localhost:8080/tasks/update


Get logs:
curl http://localhost:8080/logs


Exit server:
curl -X POST http://localhost:8080/exit \
     -H "Content-Type: application/json" \
     -d '{"password":"SERVER_EXIT_PASSWORD"}'
