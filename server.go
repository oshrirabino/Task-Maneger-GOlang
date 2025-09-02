package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AddTaskRequest struct {
	Title string `json:"title"`
}
type IDTaskRequest struct {
	ID int `json:"id"`
}

var srv *http.Server // global reference so handler can call Shutdown
var PORT = "8080"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Task Manager!")
}

func servAddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowd", http.StatusMethodNotAllowed)
	}

	var req AddTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task, err := AddTask(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logChan <- fmt.Sprintf("ADD TASK %s SUCCESSFULY WITH ID %d", task.Title, task.ID)

	// Respond with the created task as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func servCompleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowd", http.StatusMethodNotAllowed)
	}
	var req IDTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := CompleteTask(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logChan <- fmt.Sprintf("Task %d Comleted", req.ID)

	// Respond with status ok
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func servRemoveHandlder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowd", http.StatusMethodNotAllowed)
	}
	var req IDTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := RemoveTask(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logChan <- fmt.Sprintf("TASK %d REMOVED", req.ID)

	// Respond with status ok
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func servListHandler(w http.ResponseWriter, r *http.Request) {
	pending, completed := sortedTasks()
	query := r.URL.Query()

	var result []Task
	if query.Has("pending") {
		result = pending
	} else if query.Has("completed") {
		result = completed
	} else {
		// default: return both (pending first, then completed)
		result = append(pending, completed...)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func servLogsHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read logs
	logData := ReadLogs()
	if logData[:4] != "LOGS" {
		http.Error(w, logData, http.StatusInternalServerError)
		return
	}

	// Send logs as JSON
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"logs": logData}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode logs: "+err.Error(), http.StatusInternalServerError)
	}
}

func servUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Remove completed tasks
	removedTasks := RemoveComplitedTasks() // returns []Task

	// Send removed tasks as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(removedTasks); err != nil {
		http.Error(w, "Failed to encode removed tasks: "+err.Error(), http.StatusInternalServerError)
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	// Decode JSON from client
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check password
	if req.Password != PASSWORD {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Shutdown server in a goroutine
	go func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	// Respond immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"server shutting down"}`))
}

func runServer() {

	mux := http.NewServeMux()
	srv = &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/tasks/add", servAddTaskHandler)
	mux.HandleFunc("/tasks/complte", servCompleteHandler)
	mux.HandleFunc("/tasks/remove", servRemoveHandlder)
	mux.HandleFunc("/tasks/list", servListHandler)
	mux.HandleFunc("/tasks/update", servUpdateHandler)
	mux.HandleFunc("/logs", servLogsHandler)
	mux.HandleFunc("/exit", exitHandler)
	fmt.Println("Server running on http://localhost:" + PORT)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server faild: %v", err)
	}
}
