package storage

import (
	"encoding/json"
	"io"
	"os"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

// BackupRestorer different types of persistent storages
type BackupRestorer interface {
	Backup(map[string]metrics.Gauge, map[string]metrics.Counter) error
	Restore(map[string]metrics.Gauge, map[string]metrics.Counter) error
}

var _ BackupRestorer = (*FileBackuper)(nil)

// FileBackuper file backuper struct
type FileBackuper struct {
	file *os.File
}

// NewFileBackuper creates new FileBackuper instance
func NewFileBackuper(f *os.File) FileBackuper {
	return FileBackuper{
		file: f,
	}
}

// Backup put metrics to file
func (fb FileBackuper) Backup(gauges map[string]metrics.Gauge, counters map[string]metrics.Counter) error {
	var structs []metrics.Metrics
	// clear file content
	if _, err := fb.file.Seek(0, 0); err != nil {
		logger.Info("seek file error")
	}
	if err := fb.file.Truncate(0); err != nil {
		logger.Info("truncate backup file error")
		return err
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

	_, err = fb.file.WriteString(string(bytes))
	if err != nil {
		return err
	}

	logger.Info("backuped", zap.Int("values count", len(structs)))
	return nil
}

// Restore get metrics from file backup
func (fb FileBackuper) Restore(gauges map[string]metrics.Gauge, counters map[string]metrics.Counter) error {
	var structs []metrics.Metrics
	decoder := json.NewDecoder(fb.file)
	if err := decoder.Decode(&structs); err != nil && err != io.EOF {
		return err
	}
	for _, v := range structs {
		switch v.MType {
		case "gauge":
			gauges[v.ID] = metrics.Gauge(*v.Value)
		case "counter":
			counters[v.ID] = metrics.Counter(*v.Delta)
		}
	}
	logger.Info("success restored", zap.Int("metrics count", len(structs)))
	return nil
}
