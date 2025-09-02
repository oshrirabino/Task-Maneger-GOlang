package main

import (
	"errors"
	"strconv"
	"sync"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var taskFile string

var tasks []Task
var tasksMutex sync.Mutex

func nextID() int {
	maxId := 0
	for _, t := range tasks {
		if maxId < t.ID {
			maxId = t.ID
		}
	}
	return maxId + 1
}

func AddTask(title string) (Task, error) {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	if title == "" {
		return Task{}, errors.New("title cannot be empty")
	}
	task := Task{
		ID:        nextID(),
		Title:     title,
		Completed: false,
	}
	tasks = append(tasks, task)
	return task, nil
}

func ListTasks() []Task {
	return tasks
}

func CompleteTask(ID int) error {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	for i, task := range tasks {
		if task.ID == ID {
			tasks[i].Completed = true
			return nil
		}
	}
	return errors.New("task " + strconv.Itoa(ID) + " not found")
}

func RemoveTask(ID int) error {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	for i, task := range tasks {
		if task.ID == ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("Task with ID: " + strconv.Itoa(ID) + " not found")
}

func sortedTasks() ([]Task, []Task) { // first pending, next complited
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	unCompTasks := []Task{}
	compTasks := []Task{}
	for _, task := range tasks {
		if task.Completed {
			compTasks = append(compTasks, task)
		} else {
			unCompTasks = append(unCompTasks, task)
		}
	}
	return unCompTasks, compTasks
}

func RemoveComplitedTasks() []Task {
	updatedTasks, comletedTasks := sortedTasks()

	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	tasks = updatedTasks
	return comletedTasks
}

func ClearTasks() {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	tasks = []Task{}
}
