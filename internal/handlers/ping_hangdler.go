package handlers

import (
	"net/http"

	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

// Ping is a /ping/ handler, DB connect healthcheck
func (bHandler baseHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bHandler.db == nil {
			logger.Warn("db is not initialized")
			w.WriteHeader(http.StatusTeapot)
			return
		}
		_, err := bHandler.db.Exec("SELECT 1")
		if err != nil {
			logger.Error("can't ping db", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("pong"))
		if err != nil {
			logger.Error("can't write response", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
