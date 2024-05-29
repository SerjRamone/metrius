package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
	"github.com/SerjRamone/metrius/pkg/retry"
)

var _ Storage = (*SQLStorage)(nil)

// SQLStorage is a database storage
type SQLStorage struct {
	db *sql.DB
}

// NewSQLStorage creates SQL db storage
func NewSQLStorage(dsn string) (SQLStorage, error) {
	var stor SQLStorage
	if dsn == "" {
		return stor, errors.New("db dsn not provided")
	}

	pgDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return stor, err
	}

	// run test select query to make sure PostgreSQL is up and running
	err = retry.WithBackoff(func() error {
		return pgDB.Ping()
	}, 3)

	if err != nil {
		return stor, err
	}
	logger.Info("db ping attempt", zap.String("status", "OK"))

	stor = SQLStorage{
		db: pgDB,
	}

	return stor, nil
}

// SetGauge insert or update metrics value of type gauge
func (dbs SQLStorage) SetGauge(ctx context.Context, name string, value metrics.Gauge) error {
	stmt, err := dbs.db.PrepareContext(ctx, "INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()")
	if err != nil {
		logger.Error("statement creating error", zap.Error(err))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(context.TODO(), name, "gauge", float64(value))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsConnectionException(pgErr.Code) {
				err = retry.WithBackoff(func() error {
					_, err = stmt.ExecContext(context.TODO(), name, "gauge", float64(value))
					return err
				}, 3)
			}
		}
		logger.Error("db upsert error", zap.String("name", name), zap.Float64("value", float64(value)), zap.Error(err))
		return err
	}
	logger.Info("db upsert success", zap.String("name", name), zap.Float64("value", float64(value)))
	return nil
}

// Gauge returns value of type gauge by name
func (dbs SQLStorage) Gauge(ctx context.Context, name string) (metrics.Gauge, bool) {
	stmt, err := dbs.db.PrepareContext(ctx, "SELECT value FROM metrics WHERE mtype='gauge' AND id=$1")
	if err != nil {
		logger.Error("statement creating error", zap.Error(err))
		return 0, false
	}
	defer stmt.Close()
	var row *sql.Row
	row = stmt.QueryRowContext(context.TODO(), name)
	var value float64
	err = row.Scan(&value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsConnectionException(pgErr.Code) {
				err = retry.WithBackoff(func() error {
					row = stmt.QueryRowContext(context.TODO(), name)
					err = row.Scan(&value)
					if !errors.Is(err, sql.ErrNoRows) {
						return err
					}
					return nil
				}, 3)
			}
		}
		logger.Error("row scan error", zap.String("name", name), zap.Error(err))
		return 0, false
	}
	return metrics.Gauge(value), true
}

// Gauges returns map of all setted gauges
func (dbs SQLStorage) Gauges(ctx context.Context) map[string]metrics.Gauge {
	result := map[string]metrics.Gauge{}
	rows, err := dbs.db.QueryContext(ctx, "SELECT id, delta FROM metrics WHERE mtype='gauge'")
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
func (dbs SQLStorage) SetCounter(ctx context.Context, name string, value metrics.Counter) error {
	stmt, err := dbs.db.PrepareContext(ctx, "INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + metrics.delta, updated_at = NOW()")
	if err != nil {
		logger.Error("statement creating error", zap.Error(err))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(context.TODO(), name, "counter", int64(value))
	if err != nil {
		logger.Error("db upsert error", zap.String("name", name), zap.Int64("delta", int64(value)), zap.Error(err))
		return err
	}
	logger.Info("db upsert success", zap.String("name", name), zap.Int64("delta", int64(value)))
	return nil
}

// Counter returns value of type counter by name
func (dbs SQLStorage) Counter(ctx context.Context, name string) (metrics.Counter, bool) {
	stmt, err := dbs.db.PrepareContext(ctx, "SELECT delta FROM metrics WHERE mtype='counter' AND id=$1")
	if err != nil {
		logger.Error("statement creating error", zap.Error(err))
		return 0, false
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(context.TODO(), name)
	var delta int64
	err = row.Scan(&delta)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error("scan error", zap.Error(err))
		}
		return 0, false
	}
	return metrics.Counter(delta), true
}

// Counters returns map of all setted counters
func (dbs SQLStorage) Counters(ctx context.Context) map[string]metrics.Counter {
	result := map[string]metrics.Counter{}
	rows, err := dbs.db.QueryContext(ctx, "SELECT id, delta FROM metrics WHERE mtype='counter'")
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

// BatchUpsert insert or updates metrics in batches
func (dbs SQLStorage) BatchUpsert(ctx context.Context, batch []metrics.Metrics) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		logger.Error("transaction begin error", zap.Error(err))
		return err
	}

	// gauge statement
	stmtG, err := dbs.db.PrepareContext(ctx, "INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()")
	if err != nil {
		logger.Error("gauge statement creating error", zap.Error(err))
		return err
	}
	defer stmtG.Close()

	// counter statement
	stmtC, err := dbs.db.PrepareContext(ctx, "INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + metrics.delta, updated_at = NOW()")
	if err != nil {
		logger.Error("counter statement creating error", zap.Error(err))
		return err
	}
	defer stmtC.Close()

	for _, m := range batch {
		switch m.MType {
		case "gauge":
			_, err := stmtG.ExecContext(ctx, m.ID, "gauge", float64(*m.Value))
			if err != nil {
				logger.Error("batch upsert gauge error", zap.Error(err))
				if err = tx.Rollback(); err != nil {
					logger.Error("tx rollback error", zap.Error(err))
				}
				return err
			}
		case "counter":
			_, err := stmtC.ExecContext(ctx, m.ID, "counter", int64(*m.Delta))
			if err != nil {
				logger.Error("batch upsert counter error", zap.Error(err))
				if err = tx.Rollback(); err != nil {
					logger.Error("tx rollback error", zap.Error(err))
				}
				return err
			}
		default:
			return fmt.Errorf("unknown metrics type: %v", m.MType)
		}
	}

	return tx.Commit()
}

// Ping checks connection
func (dbs SQLStorage) Ping() error {
	return dbs.db.Ping()
}

// DBClose closes db connection
func (dbs SQLStorage) DBClose() error {
	return dbs.db.Close()
}
