// Package config ...
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	agentDefaultConfig         = ""

	agentUsageServerAddress  = "address and port of metrics server"
	agentUsageReportInterval = "period of time for sending data to server in seconds"
	agentUsagePollInterval   = "period of time for collecting metrics values in seconds"
	agentUsageHashKey        = "key string for hashing function"
	agentUsageRateLimit      = "number of synchronous outgoing requests"
	agentUsageCryptoKey      = "path to the public key file"
	agentUsageConfig         = "path to config.json file"

	serverDefaultAddress         = "localhost:8080"
	serverDefaultStoreInterval   = 300
	serverDefaultFileStoragePath = "/tmp/metrics-db.json"
	serverDefaultRestore         = true
	serverDefaultDatabaseDSN     = ""
	serverDefaultHashKey         = ""
	serverDefaultCryptoKey       = ""
	serverDefaultConfig          = ""
	serverDefaultTrustedSubnet   = ""

	serverUsageAddress         = "address and port to run server"
	serverUsageStoreInterval   = "period of time for put metrics to file"
	serverUsageFileStoragePath = "path to file for store metrics"
	serverUsageRestore         = "if true then server will resotre metrics from file storage"
	serverUsageDatabaseDSN     = "data sourse string in format: \"host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable\""
	serverUsageHashKey         = "key string for hashing function"
	serverUsageCryptoKey       = "path to the private key file"
	serverUsageConfig          = "path to config.json file"
	serverUsageTrustedSubnet   = "CIDR"
)

var errTypeAssert = errors.New("type assesrtion error")

// Agent contents config for Agent
type Agent struct {
	ServerAddress  string `env:"ADDRESS" json:"address"`
	HashKey        string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	RateLimit      int    `env:"RATE_LIMIT"`
	Config         string `env:"CONFIG"`
}

// NewAgent constructor for agent config
func NewAgent() (Agent, error) {
	a := Agent{}
	a.parseFlags()
	if err := a.parseEnv(); err != nil {
		return a, err
	}
	if a.Config != "" {
		if err := a.parseFile(); err != nil {
			return a, err
		}
	}
	return a, nil
}

// parseFlags parse cli flags
func (c *Agent) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", agentDefaultServerAddress, agentUsageServerAddress)
	flag.IntVar(&c.ReportInterval, "r", agentDefaultReportInterval, agentUsageReportInterval)
	flag.IntVar(&c.PollInterval, "p", agentDefaultPollInterval, agentUsagePollInterval)
	flag.StringVar(&c.HashKey, "k", agentDefaultHashKey, agentUsageHashKey)
	flag.IntVar(&c.RateLimit, "l", agentDefaultRateLimit, agentUsageRateLimit)
	flag.StringVar(&c.CryptoKey, "crypto-key", agentDefaultCryptoKey, agentUsageCryptoKey)
	flag.StringVar(&c.CryptoKey, "c", agentDefaultConfig, agentUsageConfig)
	flag.StringVar(&c.CryptoKey, "config", agentDefaultConfig, agentUsageConfig)

	flag.Parse()
}

// parseEnv parse environtment variables
func (c *Agent) parseEnv() error {
	return env.Parse(c)
}

// parseFile parse config file if path setted
func (c *Agent) parseFile() error {
	bytes, err := os.ReadFile(c.Config)
	if err != nil {
		return fmt.Errorf("reading config file <%s> error: %w", c.Config, err)
	}
	var tmp map[string]any
	err = json.Unmarshal(bytes, &tmp)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	var ok bool
	for param, val := range tmp {
		if param == "address" && c.ServerAddress == agentDefaultServerAddress {
			c.ServerAddress, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for ServerAddress, received: %T", errTypeAssert, val)
			}
		}
		if param == "report_interval" && c.ReportInterval == agentDefaultReportInterval {
			var v string
			v, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for ReportInterval, received: %T", errTypeAssert, val)
			}
			c.ReportInterval, err = parseInterval(v)
			if err != nil {
				return fmt.Errorf("parseInterval value <%s> error: %w", v, err)
			}
		}
		if param == "poll_interval" && c.PollInterval == agentDefaultPollInterval {
			var v string
			v, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for PollInterval, received: %T", errTypeAssert, val)
			}
			c.PollInterval, err = parseInterval(v)
			if err != nil {
				return fmt.Errorf("parseInterval value <%s> error: %w", v, err)
			}
		}
		if param == "crypto_key" && c.CryptoKey == agentDefaultCryptoKey {
			c.CryptoKey, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for CryptoKey, received: %T", errTypeAssert, val)
			}
		}
	}
	return nil
}

// MarshalLogObject zapcore.ObjectMarshaler implemet for loggin agent config struct
func (c *Agent) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ServerAddress", c.ServerAddress)
	enc.AddInt("ReportInterval", c.ReportInterval)
	enc.AddInt("PollInterval", c.PollInterval)
	enc.AddString("HashKey", c.HashKey)
	enc.AddInt("RateLimit", c.RateLimit)
	enc.AddString("CryptoKey", c.CryptoKey)
	enc.AddString("Config", c.Config)
	return nil
}

// Server contents config for Server
type Server struct {
	Address         string `env:"ADDRESS" json:"address"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey         string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	Config          string `env:"CONFIG"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
}

// NewServer constructor for server config
func NewServer() (Server, error) {
	s := Server{}
	s.parseFlags()
	if err := s.parseEnv(); err != nil {
		return s, err
	}
	if s.Config != "" {
		if err := s.parseFile(); err != nil {
			return s, err
		}
	}
	return s, nil
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
	flag.StringVar(&c.Config, "c", serverDefaultConfig, serverUsageConfig)
	flag.StringVar(&c.Config, "config", serverDefaultConfig, serverUsageConfig)
	flag.StringVar(&c.TrustedSubnet, "t", serverDefaultTrustedSubnet, serverUsageTrustedSubnet)

	flag.Parse()
}

// parseEnv parse environtment variables
func (c *Server) parseEnv() error {
	return env.Parse(c)
}

// parseFile parse config file if path setted
func (c *Server) parseFile() error {
	bytes, err := os.ReadFile(c.Config)
	if err != nil {
		return fmt.Errorf("reading config file <%s> error: %w", c.Config, err)
	}
	var tmp map[string]any
	err = json.Unmarshal(bytes, &tmp)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	var ok bool
	for param, val := range tmp {
		if param == "address" && c.Address == serverDefaultAddress {
			c.Address, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for Address, received: %T", errTypeAssert, val)
			}
		}
		if param == "restore" && c.Restore { // c.Restore is true by default
			c.Restore, ok = val.(bool)
			if !ok {
				return fmt.Errorf("%w: expected type bool for Restore, received: %T", errTypeAssert, val)
			}
		}
		if param == "store_interval" && c.StoreInterval == serverDefaultStoreInterval {
			var v string
			v, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for StoreInterval, received: %T", errTypeAssert, val)
			}
			c.StoreInterval, err = parseInterval(v)
			if err != nil {
				return fmt.Errorf("parseInterval value <%s> error: %w", v, err)
			}
		}
		if param == "store_file" && c.FileStoragePath == serverDefaultFileStoragePath {
			c.FileStoragePath, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for FileStoragePath, received: %T", errTypeAssert, val)
			}
		}
		if param == "database_dsn" && c.DatabaseDSN == serverDefaultDatabaseDSN {
			c.DatabaseDSN, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for DatabaseDSN, received: %T", errTypeAssert, val)
			}
		}
		if param == "crypto_key" && c.CryptoKey == serverDefaultCryptoKey {
			c.CryptoKey, ok = val.(string)
			if !ok {
				return fmt.Errorf("%w: expected type string for CryptoKey, received: %T", errTypeAssert, val)
			}
		}
	}
	return nil
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
	enc.AddString("Config", c.Config)
	enc.AddString("TrustedSubnet", c.TrustedSubnet)
	return nil
}

// parseInterval parse string interval to int interval value
func parseInterval(v string) (int, error) {
	if strings.HasSuffix(v, "s") {
		v, _ = strings.CutSuffix(v, "s")
	}
	interval, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return interval, nil
}
