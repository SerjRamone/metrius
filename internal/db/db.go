// Package db ...
package db

import (
	"database/sql"

	"github.com/SerjRamone/metrius/pkg/logger"
	"github.com/SerjRamone/metrius/pkg/retry"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// DB ...
type DB struct {
	*sql.DB
}

// Dial ...
func Dial(dsn string) (*DB, error) {
	if dsn == "" {
		// return nil, errors.New("db dsn not provided")
		return &DB{}, nil
	}

	pgDB, err := sql.Open("pgx/v5", dsn)
	if err != nil {
		return nil, err
	}

	// run test select query to make sure PostgreSQL is up and running
	const maxRetries = 3
	err = retry.WithBackoff(func() error {
		_, err = pgDB.Exec("SELECT 1")
		return err
	}, maxRetries)

	if err != nil {
		return nil, err
	}
	logger.Info("db ping attempt", zap.String("status", "OK"))

	return &DB{pgDB}, nil
}
