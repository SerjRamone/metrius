// Package main
package main

import (
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
	mStorage := storage.New()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.Update(mStorage))
	mux.HandleFunc("/", handlers.Main(mStorage))

	return http.ListenAndServe(":8080", mux)
}
