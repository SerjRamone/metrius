// Package collector ...
package collector

import (
	"log"
	"runtime"

	"github.com/SerjRamone/metrius/internal/metrics"
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
	log.Println("🗄  metrics added. len(Collector.collections) is: ", len(c.collections))
}

// Export returns collections and clear slice
func (c *collector) Export() []metrics.Collection {
	collections := c.collections
	c.collections = make([]metrics.Collection, 5)
	return collections
}
