// Package storage ...
package storage

import (
	"errors"
	"fmt"

	"github.com/SerjRamone/metrius/internal/metrics"
)

var errorStorageNotInit = errors.New("storage is not initialized")

// Storage differet types storage interface
type Storage interface {
	SetGauge(string, metrics.Gauge) error
	Gauge(string) (metrics.Gauge, bool)
	Gauges() map[string]metrics.Gauge
	SetCounter(string, metrics.Counter) error
	Counter(string) (metrics.Counter, bool)
	Counters() map[string]metrics.Counter
	String() string
}

// MemStorage is a in-memory storage
type MemStorage struct {
	gauges   map[string]metrics.Gauge
	counters map[string]metrics.Counter
}

// New is a constructor of MemStorage storage
func New() MemStorage {
	return MemStorage{
		gauges:   map[string]metrics.Gauge{},
		counters: map[string]metrics.Counter{},
	}
}

// SetGauge insert or update metrics value of type gauge
func (s MemStorage) SetGauge(name string, value metrics.Gauge) error {
	if s.gauges == nil {
		return errorStorageNotInit
	}
	s.gauges[name] = value
	return nil
}

// Gauge returns value of type gauge by name
func (s MemStorage) Gauge(name string) (v metrics.Gauge, ok bool) {
	v, ok = s.gauges[name]
	return
}

// SetCounter increase metrics value of type counter
func (s MemStorage) SetCounter(name string, value metrics.Counter) error {
	if s.counters == nil {
		return errorStorageNotInit
	}
	s.counters[name] += value
	return nil
}

// Counter returns value of type counter by name
func (s MemStorage) Counter(name string) (v metrics.Counter, ok bool) {
	v, ok = s.counters[name]
	return
}

// Gauges returns map of all setted gauges
func (s MemStorage) Gauges() map[string]metrics.Gauge {
	return s.gauges
}

// Counters returns map of all setted gauges
func (s MemStorage) Counters() map[string]metrics.Counter {
	return s.counters
}

// String returns data as string
func (s MemStorage) String() string {
	r := "Counters: \r\n"
	r += fmt.Sprintf("%v", s.counters) + "\r\n"
	r += "Gauges: \r\n"
	r += fmt.Sprintf("%v", s.gauges) + "\r\n"
	return r
}
