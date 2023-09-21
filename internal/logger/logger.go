// Package logger ...
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is no-operation logger by default
// Init in Init() function
var Log *zap.Logger = zap.NewNop()

type (
	// stores response data
	responseData struct {
		status int
		size   int
	}

	// http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write ...
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// write response via original http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // get size
	return size, err
}

// WriteHeader ...
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// write status code via original http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // get status code
}

// Init and configure logger
func Init(level string) error {
	// set log level
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

// RequestLogger middleware for request logging
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		// get request duration
		duration := time.Since(start)

		Log.Info("handle request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	})
}
