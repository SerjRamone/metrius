// Package handlers provides HTTP request handlers used throughout the project.
package handlers

import (
	"context"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/storage"
)

type metricsStorage interface {
	SetGauge(context.Context, string, metrics.Gauge) error
	Gauge(context.Context, string) (metrics.Gauge, bool)
	Gauges(context.Context) map[string]metrics.Gauge
	SetCounter(context.Context, string, metrics.Counter) error
	Counter(context.Context, string) (metrics.Counter, bool)
	Counters(context.Context) map[string]metrics.Counter
	BatchUpsert(context.Context, []metrics.Metrics) error
}

// baseHandler base handler with storage inside
type baseHandler struct {
	storage metricsStorage
}

// NewBaseHandler creates a new instance of the base handler with the specified data storage.
func NewBaseHandler(storage storage.Storage) baseHandler {
	return baseHandler{
		storage: storage,
	}
}
