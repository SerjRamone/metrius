package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/go-chi/chi/v5"
)

// Update is a /update/ handler
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

		log.Println("📨  ", r.Method, r.URL.Path)

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
			c, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := bHandler.storage.SetCounter(mName, metrics.Counter(c)); err != nil {
				log.Fatal("can't set counter", err)
				return
			}

		case "gauge":
			g, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := bHandler.storage.SetGauge(mName, metrics.Gauge(g)); err != nil {
				log.Fatal("can't set gauge", err)
				return
			}
		}

		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Println("can't write response:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
