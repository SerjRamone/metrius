// Package metrics describe possible types of metrics
package metrics

// Gauge type is a replacement type. Every value set replace prev value
type Gauge float64

// Counter type is adds new value to the prev one new
type Counter int64

// Metrics type describes JSON-request format
type Metrics struct {
	ID    string   `json:"id"`              // name of metrics
	MType string   `json:"type"`            // "gauge" or "counter"
	Delta *int64   `json:"delta,omitempty"` // metrics value if type is counter
	Value *float64 `json:"value,omitempty"` // metrics value if type is gauge
}
