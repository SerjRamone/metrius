package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/SerjRamone/metrius/pkg/logger"
)

// Ping handles GET requests to the /ping/ address, performing a health check of the database connection.
// Possible response status codes:
//   - 500 in case of an internal service error.
//   - 418 in all other cases.
func (bHandler baseHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v, ok := bHandler.storage.(storage.SQLStorage); ok {
			err := v.Ping()
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
			return
		}

		logger.Warn("storage is not a SQLStorage")
		w.WriteHeader(http.StatusTeapot)
	}
}
