package main

import (
	"flag"
)

// address and port to run server
var (
	serverAddress  string
	reportInterval int
	pollInterval   int
)

// parseFlags registers and parses arguments
func parseFlags() {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&reportInterval, "r", 10, "period of time for sending data to server in seconds")
	flag.IntVar(&pollInterval, "p", 2, "period of time for collecting metrics values in seconds")
	flag.Parse()
}
