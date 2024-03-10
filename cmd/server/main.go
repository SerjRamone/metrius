// Package main
package main

import (
	"context"
	"errors"
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
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/SerjRamone/metrius/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		panic(err)
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

	ctx, cancel := context.WithCancel(context.Background())

	// catch signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// creating server
	server := &http.Server{
		Addr:    conf.Address,
		Handler: handlers.Router(stor, conf.HashKey),
	}

	if conf.Restore {
		if v, ok := stor.(storage.MemStorage); ok {
			if err := v.Restore(); err != nil {
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

					if err := v.Backup(); err != nil {
						logger.Error("backup error", zap.Error(err))
					}
				}
			}()
		}
	}

	go func() {
		logger.Info("starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

		// shutting down server
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutting down error", zap.Error(err))
		} else {
			logger.Info("server shut down gracefully")
		}

		// backup metrics
		if v, ok := stor.(storage.MemStorage); ok {
			if err := v.Backup(); err != nil {
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
