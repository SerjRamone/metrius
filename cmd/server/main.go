// Package main
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/config"
	"github.com/SerjRamone/metrius/internal/handlers"
	"github.com/SerjRamone/metrius/internal/server"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/SerjRamone/metrius/pkg/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printTags()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := config.NewServer()
	if err != nil {
		return err
	}

	flagLogLevel := "info"
	if err = logger.Init(flagLogLevel); err != nil {
		return err
	}

	logger.Info("loaded config", zap.Object("config", &conf))

	var backuper storage.BackupRestorer
	var backupFile *os.File
	var stor storage.Storage

	if conf.DatabaseDSN != "" { // store metrics in database
		stor, err = storage.NewSQLStorage(conf.DatabaseDSN)
		if err != nil {
			return err
		}

		// run Postgres migrations
		logger.Info("running pg migrations")
		if err = runPgMigrations(conf.DatabaseDSN); err != nil {
			return err
		}
	} else { // store metrics in memory
		backupFile, err = os.OpenFile(conf.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		// backup to file
		backuper = storage.NewFileBackuper(backupFile)

		// init MemStorage
		stor = storage.NewMemStorage(conf.StoreInterval, backuper)
	}

	var privKey []byte
	if conf.CryptoKey != "" {
		privKey, err = os.ReadFile(conf.CryptoKey)
		if err != nil {
			logger.Error("reading keyfile error", zap.Error(err))
			return err
		}
	}

	var trustedSubnet *net.IPNet
	if conf.TrustedSubnet != "" {
		_, trustedSubnet, err = net.ParseCIDR(conf.TrustedSubnet)
		if err != nil {
			logger.Error("invalid CIDR subnet string format", zap.String("subnet", conf.TrustedSubnet), zap.Error(err))
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	// catch signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var serv server.Server
	if conf.Type == "http" {
		serv = server.NewHTTPServer(conf.Address, handlers.Router(stor, conf.HashKey, privKey, trustedSubnet))
	} else if conf.Type == "grpc" {
		serv = server.NewGRPCServer(conf.Address, stor)
	} else {
		logger.Error("invalid server type", zap.String("type", conf.Type))
		cancel()
	}

	if conf.Restore {
		if v, ok := stor.(storage.MemStorage); ok {
			if err := v.Restore(ctx); err != nil {
				logger.Error("can't restore from file", zap.Error(err))
				cancel()
			}
		}
	}

	if conf.StoreInterval != 0 {
		if v, ok := stor.(storage.MemStorage); ok {
			go func() {
				logger.Info("backuper started", zap.Int("StoreInterval", conf.StoreInterval))
				ticker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
				for {
					<-ticker.C

					if err := v.Backup(ctx); err != nil {
						logger.Error("backup error", zap.Error(err))
					}
				}
			}()
		}
	}

	go func() {
		logger.Info("starting server...")
		if err := serv.Up(); err != nil && err != http.ErrServerClosed {
			logger.Error("server start error", zap.Error(err))
			cancel()
		}
	}()

	// waiting signals or context done
	select {
	case <-sigCh:
		logger.Info("shutting down")

		timeout := 3 * time.Second
		shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if err := serv.Down(shutdownCtx); err != nil {
			logger.Error("server shutting down error", zap.Error(err))
		} else {
			logger.Info("server shut down gracefully")
		}

		// backup metrics
		if v, ok := stor.(storage.MemStorage); ok {
			if err := v.Backup(shutdownCtx); err != nil {
				logger.Error("backup error", zap.Error(err))
			}
		}

		// close resources
		if backupFile != nil {
			if err := backupFile.Close(); err != nil {
				logger.Error("backup file close error", zap.Error(err))
			}
		}

		if v, ok := stor.(storage.SQLStorage); ok {
			if err := v.DBClose(); err != nil {
				logger.Error("db closing error", zap.Error(err))
			}
		}

	case <-ctx.Done():
		logger.Error("context error", zap.Error(ctx.Err()))
		return ctx.Err()
	}

	return nil
}

// runPgMigrations runs Postgres migrations
func runPgMigrations(dsn string) error {
	m, err := migrate.New(
		"file://internal/migrations",
		dsn,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no change made by migration scripts")
			return nil
		}
		return err
	}
	return nil
}

func printTags() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
