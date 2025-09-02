package main

import (
	"fmt"
	"os"
)

type Config struct {
	TaskFile string `json:"task_file"`
	LogFile  string `json:"log_file"`
	Password string `json:"exit_password"`
}

var PASSWORD string

func main() {
	svrMod := false
	var confiFile string
	args := os.Args
	if len(args) > 1 && args[1] == "server" {
		svrMod = true
		confiFile = "serverConfig.json"
	} else {
		confiFile = "cliConfig.json"
	}
	cfg, err := LoadConfig(confiFile)
	if err != nil {
		fmt.Println("Unable to load configuration: " + err.Error())
	}
	taskFile = cfg.TaskFile
	logFileName = cfg.LogFile
	PASSWORD = cfg.Password

	fmt.Println("LOAD TASKS...")
	err = HandleTasksLoading()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		tasks = []Task{}
	}
	go AutoSave()
	go logThread()

	if svrMod {
		runServer()
	} else {
		runCLI()
	}

	close(logChan)
	fmt.Println("Exiting, SAVE TASKS...")
	err = SaveTasks()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("TASKS SAVED SUCCESSFULY")
	}
}
