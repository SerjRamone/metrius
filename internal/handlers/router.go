package handlers

import (
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Router returns chi.Router
func Router(mS storage.MemStorage) chi.Router {
	r := chi.NewRouter()
	bHandler := NewBaseHandler(mS)

	r.Get("/", bHandler.List())
	r.Get("/value/{type}/{name}", bHandler.Value())
	r.Post("/update/{type}/{name}/{value}", bHandler.Update())

	return r
}
