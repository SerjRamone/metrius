// Package metrics describe possible types of metrics
package metrics

// Gauge type is a replacement type. Every value set replace prev value
type Gauge float64

// Counter type is adds new value to the prev one new
type Counter int64
