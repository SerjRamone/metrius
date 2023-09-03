// Package handlers ...
package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/storage"
)

// Update is a /update/ handler
func Update(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mValue, mType string
			ok            bool
		)

		if r.Method != http.MethodPost {
			http.Error(w, "Uh oh! Bad request (method)", http.StatusBadRequest)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Bad `content-type`)", http.StatusBadRequest)
			return
		}

		chunks := parseUpdate(r.URL.Path)
		if chunks == nil {
			http.Error(w, "Can't parse URL", http.StatusBadRequest)
			return
		}

		if _, ok = chunks["metrics"]; !ok {
			http.Error(w, "Metrics name is not set", http.StatusNotFound)
		}

		if mType, ok = chunks["type"]; !ok || (mType != "counter" && mType != "gauge") {
			http.Error(w, "Metrics type is not set or unknown", http.StatusBadRequest)
		}

		if mValue, ok = chunks["value"]; !ok {
			http.Error(w, "Metrics value is not set", http.StatusBadRequest)
		}

		switch mType {
		case "counter":
			c, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
			}

			if err := s.SetCounter(chunks["metrics"], metrics.Counter(c)); err != nil {
				log.Fatal("can't set counter", err)
				return
			}
			w.WriteHeader(http.StatusOK)

		case "gauge":
			g, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := s.SetGauge(chunks["metrics"], metrics.Gauge(g)); err != nil {
				log.Fatal("can't set gauge", err)
			}
		}
		_, _ = w.Write([]byte("OK"))
	}
}

// Main "/" page handler
func Main(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(s.String()))
	}
}

func parseUpdate(path string) map[string]string {
	regex := regexp.MustCompile(`/update/(?P<type>\w+)?/(?P<metrics>\w+)?/(?P<value>[0-9\.\-]+)?`)
	matches := regex.FindStringSubmatch(path)

	if len(matches) == 0 {
		return nil
	}

	chunks := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			chunks[name] = matches[i]
		}
	}

	return chunks
}
