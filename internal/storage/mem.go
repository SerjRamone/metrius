package storage

import (
	"errors"
	"fmt"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

var (
	errorStorageNotInit         = errors.New("storage is not initialized")
	_                   Storage = (*MemStorage)(nil)
)

// MemStorage is a in-memory storage
type MemStorage struct {
	gauges        map[string]metrics.Gauge
	counters      map[string]metrics.Counter
	storeInterval int
	backuper      BackupRestorer
}

// NewMemStorage is a constructor of MemStorage storage
func NewMemStorage(storeInterval int, backuper BackupRestorer) MemStorage {
	return MemStorage{
		gauges:        map[string]metrics.Gauge{},
		counters:      map[string]metrics.Counter{},
		storeInterval: storeInterval,
		backuper:      backuper,
	}
}

// Backup persist store for MemStorage
func (s MemStorage) Backup() error {
	// make local copies of maps
	// @todo mutex in future
	gauges := make(map[string]metrics.Gauge)
	counters := make(map[string]metrics.Counter)
	for k, v := range s.Gauges() {
		gauges[k] = v
	}
	for k, v := range s.Counters() {
		counters[k] = v
	}

	return s.backuper.Backup(gauges, counters)
}

// Restore ...
func (s MemStorage) Restore() error {
	return s.backuper.Restore(s.Gauges(), s.Counters())
}

// SetGauge insert or update metrics value of type gauge
func (s MemStorage) SetGauge(name string, value metrics.Gauge) error {
	if s.gauges == nil {
		return fmt.Errorf("%w", errorStorageNotInit)
	}
	s.gauges[name] = value
	if s.storeInterval == 0 {
		if err := s.Backup(); err != nil {
			return err
		}
	}
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
		return fmt.Errorf("%w", errorStorageNotInit)
	}
	s.counters[name] += value
	if s.storeInterval == 0 {
		if err := s.Backup(); err != nil {
			return err
		}
	}

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

// Counters returns map of all setted counters
func (s MemStorage) Counters() map[string]metrics.Counter {
	return s.counters
}

// BatchUpsert insert or updates metrics in batches
func (s MemStorage) BatchUpsert(batch []metrics.Metrics) error {
	for _, m := range batch {
		switch m.MType {
		case "gauge":
			err := s.SetGauge(m.ID, metrics.Gauge(*m.Value))
			if err != nil {
				logger.Error("can't set gauge", zap.String("id", m.ID), zap.Float64("value", *m.Value), zap.Error(err))
			}
		case "counter":
			err := s.SetCounter(m.ID, metrics.Counter(*m.Delta))
			if err != nil {
				logger.Error("can't set counter", zap.String("id", m.ID), zap.Int64("delta", *m.Delta), zap.Error(err))
			}
		default:
			return fmt.Errorf("unknown metrics type: %v", m.MType)
		}
	}

	return nil
}
