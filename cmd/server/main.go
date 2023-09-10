// Package main
package main

import (
	"log"
	"net/http"

	"github.com/SerjRamone/metrius/internal/handlers"
	"github.com/SerjRamone/metrius/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	parseConfig()
	log.Printf("server started on address: %v\n", address)
	mStorage := storage.New()

	return http.ListenAndServe(address, handlers.Router(mStorage))
}
