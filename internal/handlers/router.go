package handlers

import (
	"net"

	"github.com/SerjRamone/metrius/internal/middlewares"
	"github.com/SerjRamone/metrius/internal/storage"

	// "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// Router creates and configures a chi.Router object.
//   - s: an object satisfying the storage.Storage interface, used as a storage for metrics.
//   - hashKey: a string representing the encryption key used for signing.
//   - privKey: a slice of bytes with private key for decrypt request body
//   - trustedSubnet: an IPNet represents an IP network
func Router(s storage.Storage, hashKey string, privKey []byte, trustedSubnet *net.IPNet) chi.Router {
	r := chi.NewRouter()
	bHandler := NewBaseHandler(s)

	r.Use(middlewares.RequestLogger)
	if hashKey != "" {
		r.Use(middlewares.Signer(hashKey))
	}
	if len(privKey) > 0 {
		r.Use(middlewares.Crypto(privKey))
	}
	if trustedSubnet != nil {
		r.Use(middlewares.IPWhitelist(trustedSubnet))
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
