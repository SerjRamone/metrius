// Package config ...
package config

import (
	"flag"

	"github.com/caarlos0/env"
	"go.uber.org/zap/zapcore"
)

const (
	agentDefaultServerAddress  = "localhost:8080"
	agentDefaultReportInterval = 10
	agentDefaultPollInterval   = 2
	agentDefaultHashKey        = ""
	agentDefaultRateLimit      = 1
	agentDefaultCryptoKey      = ""

	agentUsageServerAddress  = "address and port of metrics server"
	agentUsageReportInterval = "period of time for sending data to server in seconds"
	agentUsagePollInterval   = "period of time for collecting metrics values in seconds"
	agentUsageHashKey        = "key string for hashing function"
	agentUsageRateLimit      = "number of synchronous outgoing requests"
	agentUsageCryptoKey      = "path to the public key file"

	serverDefaultAddress         = "localhost:8080"
	serverDefaultStoreInterval   = 300
	serverDefaultFileStoragePath = "/tmp/metrics-db.json"
	serverDefaultRestore         = true
	serverDefaultDatabaseDSN     = ""
	serverDefaultHashKey         = ""
	serverDefaultCryptoKey       = ""

	serverUsageAddress         = "address and port to run server"
	serverUsageStoreInterval   = "period of time for put metrics to file"
	serverUsageFileStoragePath = "path to file for store metrics"
	serverUsageRestore         = "if true then server will resotre metrics from file storage"
	serverUsageDatabaseDSN     = "data sourse string in format: \"host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable\""
	serverUsageHashKey         = "key string for hashing function"
	serverUsageCryptoKey       = "path to the private key file"
)

// Agent contents config for Agent
type Agent struct {
	ServerAddress  string `env:"ADDRESS"`
	HashKey        string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
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
	flag.StringVar(&c.HashKey, "k", agentDefaultHashKey, agentUsageHashKey)
	flag.IntVar(&c.RateLimit, "l", agentDefaultRateLimit, agentUsageRateLimit)
	flag.StringVar(&c.CryptoKey, "crypto-key", agentDefaultCryptoKey, agentUsageCryptoKey)

	flag.Parse()
}

// parseEnv parse environtment variables
func (c *Agent) parseEnv() error {
	return env.Parse(c)
}

// MarshalLogObject zapcore.ObjectMarshaler implemet for loggin agent config struct
func (c *Agent) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ServerAddress", c.ServerAddress)
	enc.AddInt("ReportInterval", c.ReportInterval)
	enc.AddInt("PollInterval", c.PollInterval)
	enc.AddString("HashKey", c.HashKey)
	enc.AddInt("RateLimit", c.RateLimit)
	enc.AddString("CryptoKey", c.CryptoKey)
	return nil
}

// Server contents config for Server
type Server struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
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
	flag.IntVar(&c.StoreInterval, "i", serverDefaultStoreInterval, serverUsageStoreInterval)
	flag.StringVar(&c.FileStoragePath, "f", serverDefaultFileStoragePath, serverUsageFileStoragePath)
	flag.BoolVar(&c.Restore, "r", serverDefaultRestore, serverUsageRestore)
	flag.StringVar(&c.DatabaseDSN, "d", serverDefaultDatabaseDSN, serverUsageDatabaseDSN)
	flag.StringVar(&c.HashKey, "k", serverDefaultHashKey, serverUsageHashKey)
	flag.StringVar(&c.CryptoKey, "crypto-key", serverDefaultCryptoKey, serverUsageCryptoKey)

	flag.Parse()
}

// parseEnv parse environtment variables
func (c *Server) parseEnv() error {
	return env.Parse(c)
}

// MarshalLogObject zapcore.ObjectMarshaler implemet for loggin server config struct
func (c *Server) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Address", c.Address)
	enc.AddInt("StoreInterval", c.StoreInterval)
	enc.AddString("FileStoragePath", c.FileStoragePath)
	enc.AddBool("Restore", c.Restore)
	enc.AddString("DatabaseDSN", c.DatabaseDSN)
	enc.AddString("HashKey", c.HashKey)
	enc.AddString("CryptoKey", c.CryptoKey)
	return nil
}
