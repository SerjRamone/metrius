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

	r.Get("/", middlewares.GzipCompressor(bHandler.List()))

	r.Post("/value/", middlewares.GzipCompressor(bHandler.ValueJSON()))
	r.Post("/update/", middlewares.GzipCompressor(bHandler.UpdateJSON()))

	r.Get("/value/{type}/{name}", bHandler.Value())
	r.Post("/update/{type}/{name}/{value}", bHandler.Update())

	return r
}
