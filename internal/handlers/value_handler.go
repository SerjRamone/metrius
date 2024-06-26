package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
)

// Value returns the current value of a metric in the response body.
// It handles GET requests to the /value/{type}/{name} address.
// It expects two GET parameters in the request: type and name, both mandatory.
// Possible HTTP status codes returned:
//   - 404 if the requested metric is not found.
//   - 400 if the type is not equal to "gauge" or "counter".
//   - 500 in case of a service error.
//   - 200 if the requested value is found.
func (bHandler baseHandler) Value() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mType  = chi.URLParam(r, "type")
			mName  = chi.URLParam(r, "name")
			mValue string
		)

		switch mType {
		case "gauge":
			v, ok := bHandler.storage.Gauge(r.Context(), mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		case "counter":
			v, ok := bHandler.storage.Counter(r.Context(), mName)
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

// ValueJSON returns the current value of a metric in the response body as JSON.
// It handles POST requests to the /value/ address.
// It expects requests with Content-Type application/json.
// Possible HTTP status codes returned:
//   - 400 in the following cases:
//     1. if Content-Type is not application/json.
//     2. if an invalid JSON is passed in the request body.
//   - 404 if the requested metric is not found.
//   - 500 in case of a service error.
//   - 200 if the requested value is found.
//
// Example request body:
//
//	{
//	    "id": "GaugeBatchZip125",
//	    "type": "gauge"
//	}
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
			logger.Info("request body reading error", zap.Error(err))
			http.Error(w, "Request body reading error", http.StatusBadRequest)
			return
		}

		var req metrics.Metrics
		if err = json.Unmarshal(body, &req); err != nil {
			logger.Info("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch req.MType {
		case "gauge":
			v, ok := bHandler.storage.Gauge(r.Context(), req.ID)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			tmpV := float64(v)
			req.Value = &tmpV
		case "counter":
			v, ok := bHandler.storage.Counter(r.Context(), req.ID)
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
			logger.Error("Metrics marshalling error", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(bytes)
		if err != nil {
			logger.Error("can't write response:", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
