package metrics

import (
	"math/rand"
	"runtime"
	"time"
)

// CollectionItem stores metrics name, type and value
type CollectionItem struct {
	Name, Type string
	Value      float64
}

// Collection of tracked metrics
type Collection []CollectionItem

// NewCollection fill and returns collection
func NewCollection(m runtime.MemStats) Collection {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := Collection{
		CollectionItem{Name: "Alloc", Type: "gauge", Value: float64(m.Alloc)},
		CollectionItem{Name: "Frees", Type: "gauge", Value: float64(m.Frees)},
		CollectionItem{Name: "BuckHashSys", Type: "gauge", Value: float64(m.BuckHashSys)},
		CollectionItem{Name: "GCCPUFraction", Type: "gauge", Value: float64(m.GCCPUFraction)},
		CollectionItem{Name: "GCSys", Type: "gauge", Value: float64(m.GCSys)},
		CollectionItem{Name: "HeapAlloc", Type: "gauge", Value: float64(m.HeapAlloc)},
		CollectionItem{Name: "HeapIdle", Type: "gauge", Value: float64(m.HeapIdle)},
		CollectionItem{Name: "HeapInuse", Type: "gauge", Value: float64(m.HeapInuse)},
		CollectionItem{Name: "HeapObjects", Type: "gauge", Value: float64(m.HeapObjects)},
		CollectionItem{Name: "HeapReleased", Type: "gauge", Value: float64(m.HeapReleased)},
		CollectionItem{Name: "HeapSys", Type: "gauge", Value: float64(m.HeapSys)},
		CollectionItem{Name: "LastGC", Type: "gauge", Value: float64(m.LastGC)},
		CollectionItem{Name: "Lookups", Type: "gauge", Value: float64(m.Lookups)},
		CollectionItem{Name: "MCacheInuse", Type: "gauge", Value: float64(m.MCacheInuse)},
		CollectionItem{Name: "MCacheSys", Type: "gauge", Value: float64(m.MCacheSys)},
		CollectionItem{Name: "MSpanInuse", Type: "gauge", Value: float64(m.MSpanInuse)},
		CollectionItem{Name: "MSpanSys", Type: "gauge", Value: float64(m.MSpanSys)},
		CollectionItem{Name: "Mallocs", Type: "gauge", Value: float64(m.Mallocs)},
		CollectionItem{Name: "NextGC", Type: "gauge", Value: float64(m.NextGC)},
		CollectionItem{Name: "NumForcedGC", Type: "gauge", Value: float64(m.NumForcedGC)},
		CollectionItem{Name: "NumGC", Type: "gauge", Value: float64(m.NumGC)},
		CollectionItem{Name: "OtherSys", Type: "gauge", Value: float64(m.OtherSys)},
		CollectionItem{Name: "PauseTotalNs", Type: "gauge", Value: float64(m.PauseTotalNs)},
		CollectionItem{Name: "StackInuse", Type: "gauge", Value: float64(m.StackInuse)},
		CollectionItem{Name: "StackSys", Type: "gauge", Value: float64(m.StackSys)},
		CollectionItem{Name: "Sys", Type: "gauge", Value: float64(m.Sys)},
		CollectionItem{Name: "TotalAlloc", Type: "gauge", Value: float64(m.TotalAlloc)},
		CollectionItem{Name: "PollCount", Type: "counter", Value: float64(1)},
		CollectionItem{Name: "RandomValue", Type: "gauge", Value: r.Float64()},
	}
	return c
}
