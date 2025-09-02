package main

import (
	"fmt"
	"os"
	"sync"
)

var logFileName string

var logMutex sync.Mutex
var logChan = make(chan string) // unbuffered

func logThread() {
	logFile, err := os.Create(logFileName) // truncates file every run
	if err != nil {
		fmt.Println("Cannot open log file:", err)
		return
	}
	for log := range logChan {
		logMutex.Lock()
		if _, err := logFile.WriteString(log + "\n"); err != nil {
			fmt.Println("Log write error:", err)
		}
		logMutex.Unlock()
	}
	defer logFile.Close()
}

func ReadLogs() string {
	logMutex.Lock()
	data, err := os.ReadFile(logFileName)
	if err != nil {
		return "Unable to read logs file: " + err.Error()
	}
	logMutex.Unlock()
	return "LOGS:\n" +
		string(data) + "\n" +
		"----------------\n" +
		"END OF LOG FILE"
}
