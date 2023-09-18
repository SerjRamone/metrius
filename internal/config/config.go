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

// Agent contents config for Agent
type Agent struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

// NewAgent constructor for agent config
func NewAgent() (Agent, error) {
	a := Agent{}
	a.parseFlags()
	return a, a.parseEnv()
}

// parseFlags parse cli flags
func (c *Agent) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", agentDefaultServerAddress, agentUsageServerAddress)
	flag.IntVar(&c.ReportInterval, "r", agentDefaultReportInterval, agentUsageReportInterval)
	flag.IntVar(&c.PollInterval, "p", agentDefaultPollInterval, agentUsagePollInterval)

	flag.Parse()
}

// parseEnv parse environtment variables
func (c *Agent) parseEnv() error {
	return env.Parse(c)
}

// Server contents config for Server
type Server struct {
	Address string `env:"ADDRESS"`
}

// NewServer constructor for server config
func NewServer() (Server, error) {
	s := Server{}
	s.parseFlags()
	return s, s.parseEnv()
}

// parseFlags parse cli flags
func (c *Server) parseFlags() {
	flag.StringVar(&c.Address, "a", serverDefaultAddress, serverUsageAddress)

	flag.Parse()
}

// parseEnv parse environtment variables
func (c *Server) parseEnv() error {
	return env.Parse(c)
}
