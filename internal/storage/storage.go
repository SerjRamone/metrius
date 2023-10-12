// Package storage ...
package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SerjRamone/metrius/internal/db"
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

var (
	_ Storage = (*MemStorage)(nil)
	_ Storage = (*SQLStorage)(nil)
)

// MemStorage is a in-memory storage
type MemStorage struct {
	gauges        map[string]metrics.Gauge
	counters      map[string]metrics.Counter
	storeInterval int
	backuper      BackupRestorer
}

// SQLStorage is a database storage
type SQLStorage struct {
	db *db.DB
}

// NewSQLStorage creates SQL db storage
func NewSQLStorage(db *db.DB) SQLStorage {
	return SQLStorage{
		db: db,
	}
}

// SetGauge insert or update metrics value of type gauge
func (dbs SQLStorage) SetGauge(name string, value metrics.Gauge) error {
	_, err := dbs.db.ExecContext(
		context.TODO(),
		"INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()",
		name,
		"gauge",
		float64(value),
	)
	if err != nil {
		logger.Error("db upsert error", zap.String("name", name), zap.Float64("value", float64(value)), zap.Error(err))
		return err
	}
	logger.Info("db upsert success", zap.String("name", name), zap.Float64("value", float64(value)))
	return nil
}

// Gauge returns value of type gauge by name
func (dbs SQLStorage) Gauge(name string) (metrics.Gauge, bool) {
	row := dbs.db.QueryRowContext(context.TODO(), "SELECT value FROM metrics WHERE mtype='gauge' AND id=$1", name)
	var value float64
	err := row.Scan(&value)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error("scan error", zap.Error(err))
		}
		return 0, false
	}
	return metrics.Gauge(value), false
}

// Gauges returns map of all setted gauges
func (dbs SQLStorage) Gauges() map[string]metrics.Gauge {
	result := map[string]metrics.Gauge{}
	rows, err := dbs.db.QueryContext(context.TODO(), "SELECT id, value FROM metrics WHERE mtype='gauge'")
	if err != nil {
		logger.Error("can't do select query")
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name  string
			value float64
		)
		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return result
		}

		result[name] = metrics.Gauge(value)
	}

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		logger.Error("rows.Next error", zap.Error(err))
	}

	return result
}

// SetCounter increase metrics value of type counter
func (dbs SQLStorage) SetCounter(name string, value metrics.Counter) error {
	_, err := dbs.db.ExecContext(
		context.TODO(),
		"INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + metrics.delta, updated_at = NOW()",
		name,
		"counter",
		int64(value),
	)
	if err != nil {
		logger.Error("db upsert error", zap.String("name", name), zap.Int64("delta", int64(value)), zap.Error(err))
		return err
	}
	logger.Info("db upsert success", zap.String("name", name), zap.Int64("delta", int64(value)))
	return nil
}

// Counter returns value of type counter by name
func (dbs SQLStorage) Counter(name string) (metrics.Counter, bool) {
	row := dbs.db.QueryRowContext(context.TODO(), "SELECT delta FROM metrics WHERE mtype='counter' AND id=$1", name)
	var delta int64
	err := row.Scan(&delta)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error("scan error", zap.Error(err))
		}
		return 0, false
	}
	return metrics.Counter(delta), false
}

// Counters returns map of all setted counters
func (dbs SQLStorage) Counters() map[string]metrics.Counter {
	result := map[string]metrics.Counter{}
	rows, err := dbs.db.QueryContext(context.TODO(), "SELECT id, delta FROM metrics WHERE mtype='counter'")
	if err != nil {
		logger.Error("can't do select query")
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name  string
			delta int64
		)
		err = rows.Scan(&name, &delta)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return result
		}

		result[name] = metrics.Counter(delta)
	}

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		logger.Error("rows.Next error", zap.Error(err))
	}

	return result
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

// Counters returns map of all setted counters
func (s MemStorage) Counters() map[string]metrics.Counter {
	return s.counters
}
