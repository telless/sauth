package utils

import (
	"log"
	"fmt"
	"os"
)

const (
	LowLogLevel   = 0
	MinorLogLevel = 1
	FatalLogLevel = 2
)

func CheckError(err error, context string, level int) {
	if err != nil {
		message := fmt.Sprintf("Error: %s. Context: %s", err.Error(), context)
		switch level {
		case LowLogLevel:
			{
				Log(message)
				break
			}
		case MinorLogLevel:
			{
				Log(message) // todo duplicate with LowLevel
				break
			}
		case FatalLogLevel:
			{
				Die(message)
				break
			}
		}
	}
}

func Log(message string) {
	fmt.Fprintln(os.Stdout, message)
	log.Print(message)
}

func Die(message string) {
	fmt.Fprintln(os.Stderr, message)
	log.Fatal(message)
}
