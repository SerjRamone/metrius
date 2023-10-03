// Package collector ...
package collector

import (
	"runtime"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
)

// collector collect and store metrics
type collector struct {
	collections []metrics.Collection
}

// New creates collector instance
func New() *collector {
	return &collector{
		collections: make([]metrics.Collection, 5),
	}
}

// Collect collects metrics
func (c *collector) Collect() {
	memStat := runtime.MemStats{}
	// getting metrics from runtime
	runtime.ReadMemStats(&memStat)

	// add metrics to collection
	collection := metrics.NewCollection(memStat)
	c.collections = append(c.collections, collection)
	logger.Info("metrics added")
}

// Export returns collections and clear slice
func (c *collector) Export() []metrics.Collection {
	collections := c.collections
	c.collections = make([]metrics.Collection, 5)
	return collections
}
