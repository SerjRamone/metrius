// Package storage ...
package storage

import (
	"github.com/SerjRamone/metrius/internal/metrics"
)

// Storage differet types storage interface
type Storage interface {
	SetGauge(string, metrics.Gauge) error
	Gauge(string) (metrics.Gauge, bool)
	Gauges() map[string]metrics.Gauge
	SetCounter(string, metrics.Counter) error
	Counter(string) (metrics.Counter, bool)
	Counters() map[string]metrics.Counter
	BatchUpsert([]metrics.Metrics) error
}
