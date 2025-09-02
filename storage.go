package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var tasksFileMutex sync.Mutex

func LoadTasks() ([]Task, error) {
	tasksFileMutex.Lock()
	data, err := os.ReadFile(taskFile)
	tasksFileMutex.Unlock()
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // If file doesn't exist, return empty slice
		}
		return nil, err
	}

	tasksMutex.Lock()
	defer tasksMutex.Unlock()
	readingErr := json.Unmarshal(data, &tasks)
	if readingErr != nil {
		return nil, readingErr
	}
	return tasks, nil
}

func HandleTasksLoading() error {
	t, err := LoadTasks()
	if err != nil {
		return err
	}
	tasks = t
	return nil
}

func SaveTasks() error {
	tasksMutex.Lock()
	jsonBytes, err := json.MarshalIndent(tasks, "", "  ")
	tasksMutex.Unlock()
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %v", err)
	}

	tasksFileMutex.Lock()
	defer tasksFileMutex.Unlock()
	err = os.WriteFile(taskFile, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to save tasks: %v", err)
	}
	return nil
}

func AutoSave() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := SaveTasks(); err != nil {
			fmt.Println("Auto-save error:", err)
		}
	}
}
