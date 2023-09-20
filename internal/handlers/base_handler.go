// Package handlers ...
package handlers

import (
	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/storage"
)

type metricsStorage interface {
	SetGauge(string, metrics.Gauge) error
	Gauge(string) (metrics.Gauge, bool)
	Gauges() map[string]metrics.Gauge
	SetCounter(string, metrics.Counter) error
	Counter(string) (metrics.Counter, bool)
	Counters() map[string]metrics.Counter
}

// baseHandler base handler with storage inside
type baseHandler struct {
	storage metricsStorage
}

// NewBaseHandler creates new baseHandler
func NewBaseHandler(storage storage.Storage) baseHandler {
	return baseHandler{storage: storage}
}
