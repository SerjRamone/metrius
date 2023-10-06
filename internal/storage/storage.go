// Package storage ...
package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
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
}

var _ Storage = (*MemStorage)(nil)

// MemStorage is a in-memory storage
type MemStorage struct {
	gauges        map[string]metrics.Gauge
	counters      map[string]metrics.Counter
	backupFile    *os.File
	storeInterval int
}

// New is a constructor of MemStorage storage
func New(storeInterval int, backupFile *os.File) MemStorage {
	return MemStorage{
		gauges:        map[string]metrics.Gauge{},
		counters:      map[string]metrics.Counter{},
		backupFile:    backupFile,
		storeInterval: storeInterval,
	}
}

// Backup ...
func (s MemStorage) Backup() error {
	var structs []metrics.Metrics
	// clear file content
	if _, err := s.backupFile.Seek(0, 0); err != nil {
		logger.Info("seek file error")
	}
	if err := s.backupFile.Truncate(0); err != nil {
		logger.Info("truncate backup file error")
		return err
	}

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

	for mName, mValue := range gauges {
		fValue := float64(mValue)
		structs = append(structs, metrics.Metrics{
			ID:    mName,
			MType: "gauge",
			Value: &fValue,
		})
	}

	for mName, mValue := range counters {
		iValue := int64(mValue)
		structs = append(structs, metrics.Metrics{
			ID:    mName,
			MType: "counter",
			Delta: &iValue,
		})
	}

	bytes, err := json.Marshal(structs)
	if err != nil {
		return err
	}

	_, err = s.backupFile.WriteString(string(bytes))
	if err != nil {
		return err
	}

	logger.Info("backuped", zap.Int("values count", len(structs)))
	return nil
}

// Restore ...
func (s MemStorage) Restore() error {
	var structs []metrics.Metrics
	decoder := json.NewDecoder(s.backupFile)
	if err := decoder.Decode(&structs); err != nil && err != io.EOF {
		return err
	}
	for _, v := range structs {
		switch v.MType {
		case "gauge":
			s.gauges[v.ID] = metrics.Gauge(*v.Value)
		case "counter":
			s.counters[v.ID] = metrics.Counter(*v.Delta)
		}
	}
	logger.Info("success restored", zap.Int("metrics count", len(structs)))
	return nil
}

// SetGauge insert or update metrics value of type gauge
func (s MemStorage) SetGauge(name string, value metrics.Gauge) error {
	if s.gauges == nil {
		return errorStorageNotInit
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
		return errorStorageNotInit
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

// Counters returns map of all setted gauges
func (s MemStorage) Counters() map[string]metrics.Counter {
	return s.counters
}
