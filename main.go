package main

import (
	"flag"
	"sauth/utils"
	"sauth/configuration"
	"os"
	"log"
)

var (
	config   configuration.Config = configuration.GetConfig()
	logFile  string               = ""
	sockFile string               = ""
)

func main() {
	err := initParams()
	if err != nil {
		panic(err)
	}
	err = setLogFile(logFile)
	utils.CheckError(err, "set log file", utils.FatalLogLevel)

	serve(sockFile)
}

func initParams() (error) {
	logFileR := flag.String("log", "/tmp/sauth.log", "log-file name (default: /tmp/sauth.log)")
	sockFileR := flag.String("sock", "/tmp/sauth.sock", "socket-file name (default: /tmp/sauth.sock)")
	flag.Parse()
	logFile = *logFileR
	sockFile = *sockFileR
	return nil
}

func setLogFile(fileName string) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}
