// Package main
package main

import (
	"log"
	"net/http"

	"github.com/SerjRamone/metrius/internal/config"
	"github.com/SerjRamone/metrius/internal/handlers"
	"github.com/SerjRamone/metrius/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	conf := config.Server{}
	conf.ParseFlags()
	err := conf.ParseEnv()
	if err != nil {
		log.Fatal("config parse error", err)
	}
	log.Printf("Loaded server config: %+v\n", conf)

	mStorage := storage.New()

	return http.ListenAndServe(conf.Address, handlers.Router(mStorage))
}
