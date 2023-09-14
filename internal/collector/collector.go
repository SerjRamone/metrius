package collector

import (
	"log"
	"runtime"
	"time"

	"github.com/SerjRamone/metrius/internal/metrics"
)

// Collector collect and store metrics
type Collector struct {
	collections map[int64]metrics.Collection
}

// NewCollector creates Collector instance
func NewCollector() *Collector {
	return &Collector{
		collections: make(map[int64]metrics.Collection),
	}
}

// Collect collects metrics
func (c Collector) Collect() {
	memStat := runtime.MemStats{}
	// getting metrics from runtime
	runtime.ReadMemStats(&memStat)

	// add metrics to collection
	collection := metrics.NewCollection(memStat)
	c.collections[time.Now().UnixMicro()] = collection
	log.Println("🗄  metrics added. len(Collector.collections) is: ", len(c.collections))
}

// Export returns collections and lear map
func (c *Collector) Export() map[int64]metrics.Collection {
	collections := c.collections
	c.collections = make(map[int64]metrics.Collection)
	return collections
}
