package handlers

import (
	"github.com/SerjRamone/metrius/internal/middlewares"
	"github.com/SerjRamone/metrius/internal/storage"

	// "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// Router returns chi.Router
func Router(s storage.Storage, hashKey string) chi.Router {
	r := chi.NewRouter()
	bHandler := NewBaseHandler(s)

	r.Use(middlewares.RequestLogger)
	if hashKey != "" {
		r.Use(middlewares.Signer(hashKey))
	}

	// r.Mount("/debug", middleware.Profiler())
  
	// with gzip compression
	r.Group(func(r chi.Router) {
		r.Use(middlewares.GzipCompressor)
		r.Get("/", bHandler.List())
		r.Post("/value/", bHandler.ValueJSON())
		r.Post("/update/", bHandler.UpdateJSON())
		r.Post("/updates/", bHandler.Updates())
	})

	r.Get("/value/{type}/{name}", bHandler.Value())
	r.Post("/update/{type}/{name}/{value}", bHandler.Update())

	r.Get("/ping", bHandler.Ping())

	return r
}
