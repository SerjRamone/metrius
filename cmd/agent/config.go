package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 10
	defaultPollInterval   = 2
)

// address and port to run server
var (
	serverAddress  string
	reportInterval int
	pollInterval   int
)

// parseConfig flags and envs
func parseConfig() {
	config := struct {
		ServerAddress  string `env:"ADDRESS"`
		ReportInterval int    `env:"REPORT_INTERVAL"`
		PollInterval   int    `env:"POLL_INTERVAL"`
	}{}

	flag.StringVar(&serverAddress, "a", defaultServerAddress, "address and port of metrics server")
	flag.IntVar(
		&reportInterval,
		"r",
		defaultReportInterval,
		"period of time for sending data to server in seconds",
	)
	flag.IntVar(
		&pollInterval,
		"p",
		defaultPollInterval,
		"period of time for collecting metrics values in seconds",
	)
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		log.Println("env parse error. will be used flags. error: ", err)
	}

	if config.ServerAddress != "" {
		serverAddress = config.ServerAddress
	}
	if config.PollInterval != 0 {
		pollInterval = config.PollInterval
	}
	if config.ReportInterval != 0 {
		reportInterval = config.ReportInterval
	}

	log.Printf(
		"serverAddress: %s, pollInterval: %v, reportInterval: %v",
		serverAddress,
		pollInterval,
		reportInterval,
	)
}
