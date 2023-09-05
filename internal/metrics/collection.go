package metrics

import (
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

// Collection of tracked metrics
type Collection []map[string]string

// NewCollection fill and returns collection
func NewCollection(m runtime.MemStats) Collection {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := Collection{
		map[string]string{
			"name":  "Alloc",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.Alloc), 'f', -1, 64),
		},
		map[string]string{
			"name":  "Frees",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.Frees), 'f', -1, 64),
		},
		map[string]string{
			"name":  "BuckHashSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.BuckHashSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "GCCPUFraction",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.GCCPUFraction), 'f', -1, 64),
		},
		map[string]string{
			"name":  "GCSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.GCSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "HeapAlloc",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.HeapAlloc), 'f', -1, 64),
		},
		map[string]string{
			"name":  "HeapIdle",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.HeapIdle), 'f', -1, 64),
		},
		map[string]string{
			"name":  "HeapInuse",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.HeapInuse), 'f', -1, 64),
		},
		map[string]string{
			"name":  "HeapObjects",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.HeapObjects), 'f', -1, 64),
		},
		map[string]string{
			"name":  "HeapReleased",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.HeapReleased), 'f', -1, 64),
		},
		map[string]string{
			"name":  "HeapSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.HeapSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "LastGC",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.LastGC), 'f', -1, 64),
		},
		map[string]string{
			"name":  "Lookups",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.Lookups), 'f', -1, 64),
		},
		map[string]string{
			"name":  "MCacheInuse",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.MCacheInuse), 'f', -1, 64),
		},
		map[string]string{
			"name":  "MCacheSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.MCacheSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "MSpanInuse",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.MSpanInuse), 'f', -1, 64),
		},
		map[string]string{
			"name":  "MSpanSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.MSpanSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "Mallocs",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.Mallocs), 'f', -1, 64),
		},
		map[string]string{
			"name":  "NextGC",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.NextGC), 'f', -1, 64),
		},
		map[string]string{
			"name":  "NumForcedGC",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.NumForcedGC), 'f', -1, 64),
		},
		map[string]string{
			"name":  "NumGC",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.NumGC), 'f', -1, 64),
		},
		map[string]string{
			"name":  "OtherSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.OtherSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "PauseTotalNs",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.PauseTotalNs), 'f', -1, 64),
		},
		map[string]string{
			"name":  "StackInuse",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.StackInuse), 'f', -1, 64),
		},
		map[string]string{
			"name":  "StackSys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.StackSys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "Sys",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.Sys), 'f', -1, 64),
		},
		map[string]string{
			"name":  "TotalAlloc",
			"type":  "gauge",
			"value": strconv.FormatFloat(float64(m.TotalAlloc), 'f', -1, 64),
		},
		map[string]string{
			"name":  "PollCount",
			"type":  "counter",
			"value": "1",
		},
		map[string]string{
			"name":  "RandomValue",
			"type":  "gauge",
			"value": strconv.FormatFloat(r.Float64(), 'f', -1, 64),
		},
	}
	return c
}
