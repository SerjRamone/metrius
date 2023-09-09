package main

import (
	"flag"
)

// address and port to run server
var address string

// parseFlags registers and parses arguments
func parseFlags() {
	flag.StringVar(&address, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
