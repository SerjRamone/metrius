// Package main
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SerjRamone/metrius/internal/config"
	"github.com/SerjRamone/metrius/internal/handlers"
	"github.com/SerjRamone/metrius/internal/logger"
	"github.com/SerjRamone/metrius/internal/storage"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	conf, err := config.NewServer()

	flagLogLevel := "info"
	if err := logger.Init(flagLogLevel); err != nil {
		return err
	}

	if err != nil {
		logger.Log.Fatal("config parse error: ", zap.Error(err))
	}

	backupFile, err := os.OpenFile(conf.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	mStorage := storage.New(conf.StoreInterval, backupFile)

	ctx, cancel := context.WithCancel(context.Background())

	// catch signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// creating server
	server := &http.Server{
		Addr:    conf.Address,
		Handler: handlers.Router(mStorage),
	}

	if conf.Restore {
		if err := mStorage.Restore(); err != nil {
			logger.Log.Error("can't restore from file", zap.Error(err))
			cancel()
		}
	}

	if conf.StoreInterval != 0 {
		go func() {
			logger.Log.Info("backuper started", zap.Int("StoreInterval", conf.StoreInterval))
			storedAt := time.Now()
			for {
				seconds := int((time.Since(storedAt)).Seconds())
				if seconds >= conf.StoreInterval {
					if err := mStorage.Backup(); err != nil {
						logger.Log.Error("backup error", zap.Error(err))
					}
					storedAt = time.Now()
				}

				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	go func() {
		logger.Log.Info("starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server start error", zap.Error(err))
			cancel()
		}
	}()

	// waiting signals or context done
	select {
	case <-sigCh:
		logger.Log.Info("shutting down")

		timeout := 3 * time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// shutting down server
		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Error("server shutting down error", zap.Error(err))
		} else {
			logger.Log.Info("server shut down gracefully")
		}

		// backup metrics
		if err := mStorage.Backup(); err != nil {
			logger.Log.Error("backup error", zap.Error(err))
		}

		// close resources
		if err := backupFile.Close(); err != nil {
			logger.Log.Error("backup file close error", zap.Error(err))
		}

	case <-ctx.Done():
		logger.Log.Error("context error", zap.Error(ctx.Err()))
		return err
	}

	return nil
}
