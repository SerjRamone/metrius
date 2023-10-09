package db

import (
	"database/sql"
	"time"

	"github.com/SerjRamone/metrius/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DB struct {
	*sql.DB
}

// Dial ...
func Dial(dsn string) (*DB, error) {
	if dsn == "" {
		// return nil, errors.New("db dsn not provided")
		return &DB{}, nil
	}

	pgDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// run test select query to make sure PostgreSQL is up and running
	var attempt uint
	const maxAttempts = 5

	for {
		attempt++

		logger.Info("db ping attempt", zap.Uint("attempt", attempt))

		_, err = pgDB.Exec("SELECT 1")
		if err != nil {
			logger.Warn("db ping attempt", zap.Uint("attempt", attempt), zap.Error(err))

			if attempt < maxAttempts {
				time.Sleep(1 * time.Second)

				continue
			}

			return nil, err
		}

		logger.Info("db ping attempt", zap.Uint("attempt", attempt), zap.String("status", "OK"))

		break
	}

	return &DB{pgDB}, nil
}
