package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"net"
	"fmt"
)

var config baseConfig = getConfig()

func main() {
	sockFile := flag.String("sock", "/tmp/sauth.sock", "socket-file name (default: /tmp/sauth.sock)")
	flag.Parse()

	os.Remove(*sockFile)
	listener, err := net.Listen("unix", *sockFile)
	if err != nil {
		panic(err)
	}
	logMessage(fmt.Sprintf("Server started on %s", *sockFile))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		ln.Close()
		os.Remove(*sockFile)
		logNotify(fmt.Sprintf("Caught signal %s: shutting down.", sig))
		os.Exit(0)
	}(listener, signals)

	startServer(listener)
}
