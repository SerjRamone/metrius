package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/SerjRamone/metrius/internal/logger"
	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Value handler handls GET "/value/counter/foo" requests
func (bHandler baseHandler) Value() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mType  = chi.URLParam(r, "type")
			mName  = chi.URLParam(r, "name")
			mValue string
		)

		switch mType {
		case "gauge":
			v, ok := bHandler.storage.Gauge(mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		case "counter":
			v, ok := bHandler.storage.Counter(mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		default:
			http.Error(w, "unknown type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mValue))
		if err != nil {
			log.Println("can't write response:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

// Value handler handls GET "/value" requests with JSON-body
func (bHandler baseHandler) ValueJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Bad content-type", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Info("request body reading error", zap.Error(err))
			http.Error(w, "Request body reading error", http.StatusBadRequest)
			return
		}

		var req metrics.Metrics
		if err := json.Unmarshal(body, &req); err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch req.MType {
		case "gauge":
			v, ok := bHandler.storage.Gauge(req.ID)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			tmpV := float64(v)
			req.Value = &tmpV
		case "counter":
			v, ok := bHandler.storage.Counter(req.ID)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			tmpV := int64(v)
			req.Delta = &tmpV
		default:
			http.Error(w, "unknown type", http.StatusBadRequest)
			return
		}

		bytes, err := json.Marshal(req)
		if err != nil {
			logger.Log.Error("Metrics marshalling error", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(bytes)
		if err != nil {
			logger.Log.Error("can't write response:", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
