package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func HandleAdd(title string) string {
	task, err := AddTask(title)
	if err != nil {
		return "ADD ERROR: " + err.Error()
	} else {
		msg := fmt.Sprintf("ADD TASK %s SUCCESSFULY WITH ID %d", task.Title, task.ID)
		logChan <- msg
		return msg
	}
}

func HandleList(f string) string {
	log := "List Of Tasks:"
	var list []Task
	switch f {
	case "--pending":
		list, _ = sortedTasks()
	case "--completed":
		_, list = sortedTasks()
	default:
		list = ListTasks()
	}
	for i, t := range list {
		log += fmt.Sprintf("\n %d. %s [%d]. (completed=%v)", i+1, t.Title, t.ID, t.Completed)
	}
	logChan <- "LIST TASKS"
	return log
}

func HandleComplete(ID int) string {
	err := CompleteTask(ID)
	if err == nil {
		msg := fmt.Sprintf("TASK %d COMPLETED", ID)
		logChan <- msg
		return msg
	}
	return "COMPLETE ERROR: " + err.Error()
}

func HandleRemove(ID int) string {
	err := RemoveTask(ID)
	if err != nil {
		return err.Error()
	}
	logChan <- fmt.Sprintf("TASK %d REMOVED", ID)
	return "Task removed"
}

func HandleUpdate() string {
	complitedTasks := RemoveComplitedTasks()
	var msg string
	if len(complitedTasks) > 0 {
		msg = "Removed Complited Tasks:"
		for _, task := range complitedTasks {
			msg += ("\n " + task.Title)
		}
	} else {
		msg = "No tasks completed"
	}
	logChan <- "LIST UPDATED"
	return msg
}

func HandleCommand(args []string) string {
	cmd := args[0]
	var msg string
	switch cmd {
	case "help", "h":
		msg = "Task Manager CLI\n\n" +
			"Usage:\n" +
			"  add <task title>                  -Add a new task\n" +
			"  remove <task id>                  -Remove exist task\n" +
			"  list <flags: --complited/pending> -List all tasks\n" +
			"  complete <task ID>                -Mark a task as completed\n" +
			"  update                            -Remove all completed tasks\n" +
			"  log(s)                            -Print all logs from the log file\n" +
			"  exit/quit/q                       -Quit the program and save current data\n" +
			"  help                              -Show this help message"
		// CLEAR option not here cause its a secret for the people who read the code =]

	case "add":
		taskName := strings.Join(args[1:], " ")
		msg = HandleAdd(taskName)

	case "list":
		var f string
		idx := FlagIndex(args)
		if idx >= 0 {
			f = args[idx]
		} else {
			f = ""
		}
		msg = HandleList(f)

	case "complete":
		if len(args) >= 2 {
			id, err := strconv.Atoi(args[1])
			if err == nil {
				msg = HandleComplete(id)
			} else {
				// Atoi faild
				msg = "Value Error: second argument must be a number. argument: " + args[1] + " error: " + err.Error()
			}
		} else {
			msg = "Wrong number of arguments. Usage: complete <ID: number>"
		}

	case "remove":
		if len(args) >= 2 {
			id, err := strconv.Atoi(args[1])
			if err == nil {
				msg = HandleRemove(id)
			} else {
				// Atoi faild
				msg = "Value Error: second argument must be a number. argument: " + args[1] + " error: " + err.Error()
			}
		} else {
			msg = "Wrong number of arguments. Usage: remove <ID: number>"
		}

	case "update":
		msg = HandleUpdate()

	case "log", "logs":
		msg = ReadLogs()

	case "quit", "exit", "q":
		msg = ""
	case "CLEAR":
		if len(args) < 2 || args[1] != PASSWORD {
			msg = "Unauthorized"
		} else {
			ClearTasks()
			msg = "Cleared all tasks. You are free"
		}
	default:
		msg = "Usage: <command> <args>\n" +
			"For help use command 'help'"
	}
	return msg
}

func runCLI() {

	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		args := strings.Fields(line)
		msg := HandleCommand(args)
		if msg != "" {
			fmt.Println(msg)
		} else {
			break
		}
	}
}
