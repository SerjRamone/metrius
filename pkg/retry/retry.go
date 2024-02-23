// Package retry ...
package retry

import (
	"time"

	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

// WithBackoff ...
func WithBackoff(retryable func() error, maxRetries int) error {
	for attempt := 1; ; attempt++ {
		err := retryable()
		if err == nil {
			return nil
		}

		logger.Error("retryable attempts", zap.Int("attempt num", attempt), zap.Error(err))

		if attempt >= maxRetries {
			return err
		}

		// increase interval before next attempt
		delay := time.Duration(attempt*(attempt+1)/2) * time.Second
		time.Sleep(delay)
	}
}
