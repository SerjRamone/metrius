package main

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

const defaultAddress = "localhost:8080"

// address and port to run the server
var address string

// parse flags and envs
func parse() {
	config := struct {
		Address string `env:"ADDRESS"`
	}{}

	flag.StringVar(&address, "a", defaultAddress, "address and port to run server")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		log.Println("env parse error. will be used flags. error: ", err)
	}
	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		address = envAddr
	}
}
