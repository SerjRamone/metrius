package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
)

// Update handles POST requests to the /update/{type}/{name}/{value} address, updating the value of a metric.
// It expects three parameters in the URL: value, type, and name, all of which are mandatory.
// Possible HTTP status codes returned:
//   - 404 if the name parameter is not provided.
//   - 400 in the following cases:
//     1. if type is not equal to "gauge" or "counter".
//     2. if the value parameter is not provided.
//     3. if the value parameter cannot be converted to a valid value of counter or gauge types.
//   - 200 in case of successful metric update.
//   - 500 in case of a service error.
func (bHandler baseHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mValue = chi.URLParam(r, "value")
			mType  = chi.URLParam(r, "type")
			mName  = chi.URLParam(r, "name")
		)

		// @todo
		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	http.Error(w, "Bad content-type", http.StatusBadRequest)
		// 	return
		// }

		if mName == "" {
			http.Error(w, "Metrics name not set", http.StatusNotFound)
			return
		}

		if mType != "counter" && mType != "gauge" {
			http.Error(w, "Metrics type not set or unknown", http.StatusBadRequest)
			return
		}

		if mValue == "" {
			http.Error(w, "Metrics value not set or invalid", http.StatusBadRequest)
			return
		}

		switch mType {
		case "counter":
			c, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := bHandler.storage.SetCounter(r.Context(), mName, metrics.Counter(c)); err != nil {
				log.Fatal("can't set counter", err)
				return
			}

		case "gauge":
			g, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := bHandler.storage.SetGauge(r.Context(), mName, metrics.Gauge(g)); err != nil {
				log.Fatal("can't set gauge", err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Println("can't write response:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

// Updates handles POST requests to the /updates/ address, allowing bulk updates of metric values in batches.
// It expects requests with Content-Type application/json.
// Possible HTTP status codes returned:
//   - 400 in the following cases:
//     1. if Content-Type is not application/json.
//     2. if an invalid JSON is passed in the request body.
//     3. if id, type, and delta/value are not correctly specified in the request body.
//   - 200 in case of successful metric update.
//
// Example request body:
// [
//
//	{
//	    "id": "CounterBatchZip48",
//	    "type": "counter",
//	    "delta": 1495722356
//	},
//	{
//	    "id": "GaugeBatchZip125",
//	    "type": "gauge",
//	    "value": 322489.8871640195
//	}
//
// ]
func (bHandler baseHandler) Updates() http.HandlerFunc {
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

		var batch []metrics.Metrics

		if err := json.Unmarshal(body, &batch); err != nil {
			logger.Info("cannot decode request JSON body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := bHandler.storage.BatchUpsert(r.Context(), batch); err != nil {
			logger.Info("cannot do batch upsert", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// UpdateJSON handles POST requests to the /update/ address, updating the value of a metric.
// It expects requests with Content-Type application/json.
// Possible HTTP status codes returned:
//   - 404 if the metric identifier is not provided.
//   - 400 in the following cases:
//     1. if Content-Type is not application/json.
//     2. if an invalid JSON is passed in the request body.
//     3. if id, type, and delta/value are not correctly specified in the request body.
//   - 500 in case of an internal service error.
//   - 200 in case of a successful metric update.
//
// Example request body:
//
//	{
//	    "id": "foo12",
//	    "type": "gauge",
//	    "value": 1.001
//	}
func (bHandler baseHandler) UpdateJSON() http.HandlerFunc {
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

		if req.ID == "" {
			http.Error(w, "Metrics name not set", http.StatusNotFound)
			return
		}

		if req.MType != "counter" && req.MType != "gauge" {
			http.Error(w, "Metrics type not set or unknown", http.StatusBadRequest)
			return
		}

		if (req.MType == "gauge" && req.Value == nil) || (req.MType == "counter" && req.Delta == nil) {
			http.Error(w, "Metrics value not set or invalid", http.StatusBadRequest)
			return
		}

		switch req.MType {
		case "counter":
			if err = bHandler.storage.SetCounter(r.Context(), req.ID, metrics.Counter(*req.Delta)); err != nil {
				logger.Fatal("can't set counter", zap.Error(err))
				return
			}
			// set new value for response
			newValue, ok := bHandler.storage.Counter(r.Context(), req.ID)
			if !ok {
				logger.Info("can't get new value of counter",
					zap.String("req.ID", req.ID),
					zap.Int64("rec.Delta", *req.Delta),
				)
			}
			intValue := int64(newValue)
			req.Delta = &intValue

		case "gauge":
			if err = bHandler.storage.SetGauge(r.Context(), req.ID, metrics.Gauge(*req.Value)); err != nil {
				logger.Fatal("can't set gauge", zap.Error(err))
				return
			}
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
