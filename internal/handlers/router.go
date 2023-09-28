package handlers

import (
	"github.com/SerjRamone/metrius/internal/middlewares"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Router returns chi.Router
func Router(mS storage.MemStorage) chi.Router {
	r := chi.NewRouter()
	bHandler := NewBaseHandler(mS)

	r.Use(middlewares.RequestLogger)
	r.Use(middlewares.GzipCompressor)

	r.Get("/", bHandler.List())

	r.Post("/value/", bHandler.ValueJSON())
	r.Get("/value/{type}/{name}", bHandler.Value())

	r.Post("/update/", bHandler.UpdateJSON())
	r.Post("/update/{type}/{name}/{value}", bHandler.Update())

	return r
}
