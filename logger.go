package main

import (
	"fmt"
	"os"
	"time"
)

const (
	lowLogLevel       = 0
	minorLogLevel     = 1
	fatalLogLevel     = 2
	projectTimeFormat = "2006-01-02 15:04:05"
)

func checkError(err error, context string, level int) {
	if err != nil {
		message := fmt.Sprintf("error: %s | context: %s", err.Error(), context)
		switch level {
		case lowLogLevel:
			logMessage(message)
		case minorLogLevel:
			logNotify(message)
		case fatalLogLevel:
			logError(message)
		}
	}
}

func logMessage(message string) {
	fmt.Fprintln(os.Stdout, currentTimeAsString(), message)
}

func logNotify(message string) {
	fmt.Fprintln(os.Stderr, currentTimeAsString(), message)
}

func logError(message string) {
	fmt.Fprintln(os.Stderr, currentTimeAsString(), message)
}

func currentTimeAsString() string {
	return time.Now().Format(projectTimeFormat)
}
