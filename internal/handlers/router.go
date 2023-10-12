package handlers

import (
	"github.com/SerjRamone/metrius/internal/db"
	"github.com/SerjRamone/metrius/internal/middlewares"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Router returns chi.Router
func Router(s storage.Storage, db *db.DB) chi.Router {
	r := chi.NewRouter()
	bHandler := NewBaseHandler(s, db)

	r.Use(middlewares.RequestLogger)

	r.Get("/", middlewares.GzipCompressor(bHandler.List()))

	r.Post("/value/", middlewares.GzipCompressor(bHandler.ValueJSON()))
	r.Post("/update/", middlewares.GzipCompressor(bHandler.UpdateJSON()))

	r.Get("/value/{type}/{name}", bHandler.Value())
	r.Post("/update/{type}/{name}/{value}", bHandler.Update())

	r.Get("/ping", bHandler.Ping())

	return r
}
