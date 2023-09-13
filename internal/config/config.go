// Package config ...
package config

import (
	"flag"

	"github.com/caarlos0/env"
)

const (
	agentDefaultServerAddress  = "localhost:8080"
	agentDefaultReportInterval = 10
	agentDefaultPollInterval   = 2

	agentUsageServerAddress  = "address and port of metrics server"
	agentUsageReportInterval = "period of time for sending data to server in seconds"
	agentUsagePollInterval   = "period of time for collecting metrics values in seconds"

	serverDefaultAddress = "localhost:8080"

	serverUsageAddress = "address and port to run server"
)

// Agent contents config for agent
type Agent struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

// ParseFlags parse cli flags
func (c *Agent) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", agentDefaultServerAddress, agentUsageServerAddress)
	flag.IntVar(&c.ReportInterval, "r", agentDefaultReportInterval, agentUsageReportInterval)
	flag.IntVar(&c.PollInterval, "p", agentDefaultPollInterval, agentUsagePollInterval)

	flag.Parse()
}

// ParseEnv parse environtment variables
func (c *Agent) ParseEnv() error {
	return env.Parse(c)
}

// Server contents config for server
type Server struct {
	Address string `env:"ADDRESS"`
}

// ParseFlags parse cli flags
func (c *Server) ParseFlags() {
	flag.StringVar(&c.Address, "a", serverDefaultAddress, serverUsageAddress)

	flag.Parse()
}

// ParseEnv parse environtment variables
func (c *Server) ParseEnv() error {
	return env.Parse(c)
}
