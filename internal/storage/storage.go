// Package storage ...
package storage

import (
	"context"

	"github.com/SerjRamone/metrius/internal/metrics"
)

// Storage differet types storage interface
type Storage interface {
	SetGauge(context.Context, string, metrics.Gauge) error
	Gauge(context.Context, string) (metrics.Gauge, bool)
	Gauges(context.Context) map[string]metrics.Gauge
	SetCounter(context.Context, string, metrics.Counter) error
	Counter(context.Context, string) (metrics.Counter, bool)
	Counters(context.Context) map[string]metrics.Counter
	BatchUpsert(context.Context, []metrics.Metrics) error
}
