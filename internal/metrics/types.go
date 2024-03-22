// Package metrics describe possible types of metrics
package metrics

import "encoding/json"

// Gauge type is a replacement type. Every value set replace prev value
type Gauge float64

// Counter type is adds new value to the prev one new
type Counter int64

// Metrics type describes JSON-request format
type Metrics struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func (m Metrics) MarshalJSON() ([]byte, error) {
	type Alias Metrics
	aux := struct {
		*Alias
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}{
		Alias: (*Alias)(&m),
		Delta: m.Delta,
		Value: m.Value,
	}
	return json.Marshal(&aux)
}
