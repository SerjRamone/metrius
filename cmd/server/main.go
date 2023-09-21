// Package main
package main

import (
	"net/http"

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

	mStorage := storage.New()

	return http.ListenAndServe(conf.Address, handlers.Router(mStorage))
}
