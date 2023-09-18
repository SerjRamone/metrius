package metrics

import (
	"math/rand"
	"runtime"
	"time"
)

// CollectionItem stores metrics name, type and value
type CollectionItem struct {
	Name, Variation string
	Value           float64
}

// Collection of tracked metrics
type Collection []CollectionItem

// NewCollection fill and returns collection
func NewCollection(m runtime.MemStats) Collection {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := Collection{
		CollectionItem{Name: "Alloc", Variation: "gauge", Value: float64(m.Alloc)},
		CollectionItem{Name: "Frees", Variation: "gauge", Value: float64(m.Frees)},
		CollectionItem{Name: "BuckHashSys", Variation: "gauge", Value: float64(m.BuckHashSys)},
		CollectionItem{Name: "GCCPUFraction", Variation: "gauge", Value: float64(m.GCCPUFraction)},
		CollectionItem{Name: "GCSys", Variation: "gauge", Value: float64(m.GCSys)},
		CollectionItem{Name: "HeapAlloc", Variation: "gauge", Value: float64(m.HeapAlloc)},
		CollectionItem{Name: "HeapIdle", Variation: "gauge", Value: float64(m.HeapIdle)},
		CollectionItem{Name: "HeapInuse", Variation: "gauge", Value: float64(m.HeapInuse)},
		CollectionItem{Name: "HeapObjects", Variation: "gauge", Value: float64(m.HeapObjects)},
		CollectionItem{Name: "HeapReleased", Variation: "gauge", Value: float64(m.HeapReleased)},
		CollectionItem{Name: "HeapSys", Variation: "gauge", Value: float64(m.HeapSys)},
		CollectionItem{Name: "LastGC", Variation: "gauge", Value: float64(m.LastGC)},
		CollectionItem{Name: "Lookups", Variation: "gauge", Value: float64(m.Lookups)},
		CollectionItem{Name: "MCacheInuse", Variation: "gauge", Value: float64(m.MCacheInuse)},
		CollectionItem{Name: "MCacheSys", Variation: "gauge", Value: float64(m.MCacheSys)},
		CollectionItem{Name: "MSpanInuse", Variation: "gauge", Value: float64(m.MSpanInuse)},
		CollectionItem{Name: "MSpanSys", Variation: "gauge", Value: float64(m.MSpanSys)},
		CollectionItem{Name: "Mallocs", Variation: "gauge", Value: float64(m.Mallocs)},
		CollectionItem{Name: "NextGC", Variation: "gauge", Value: float64(m.NextGC)},
		CollectionItem{Name: "NumForcedGC", Variation: "gauge", Value: float64(m.NumForcedGC)},
		CollectionItem{Name: "NumGC", Variation: "gauge", Value: float64(m.NumGC)},
		CollectionItem{Name: "OtherSys", Variation: "gauge", Value: float64(m.OtherSys)},
		CollectionItem{Name: "PauseTotalNs", Variation: "gauge", Value: float64(m.PauseTotalNs)},
		CollectionItem{Name: "StackInuse", Variation: "gauge", Value: float64(m.StackInuse)},
		CollectionItem{Name: "StackSys", Variation: "gauge", Value: float64(m.StackSys)},
		CollectionItem{Name: "Sys", Variation: "gauge", Value: float64(m.Sys)},
		CollectionItem{Name: "TotalAlloc", Variation: "gauge", Value: float64(m.TotalAlloc)},
		CollectionItem{Name: "PollCount", Variation: "gauge", Value: float64(1)},
		CollectionItem{Name: "RandomValue", Variation: "gauge", Value: r.Float64()},
	}
	return c
}
