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
			mValue, mType, mName string
			ok                   bool
		)

		if r.Method != http.MethodPost {
			http.Error(w, "Bad method", http.StatusBadRequest)
			return
		}

		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	http.Error(w, "Bad content-type", http.StatusBadRequest)
		// 	return
		// }

		log.Println("📨  ", r.URL.Path)
		chunks := parseUpdate(r.URL.Path)
		if chunks == nil {
			http.Error(w, "Can't parse URL", http.StatusBadRequest)
			return
		}

		if mName, ok = chunks["metrics"]; !ok || mName == "" {
			http.Error(w, "Metrics name not set", http.StatusNotFound)
			return
		}

		if mType, ok = chunks["type"]; !ok || (mType != "counter" && mType != "gauge") {
			http.Error(w, "Metrics type not set or unknown", http.StatusBadRequest)
			return
		}

		if mValue, ok = chunks["value"]; !ok || mValue == "" {
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

			if err := s.SetCounter(chunks["metrics"], metrics.Counter(c)); err != nil {
				log.Fatal("can't set counter", err)
				return
			}

		case "gauge":
			g, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := s.SetGauge(chunks["metrics"], metrics.Gauge(g)); err != nil {
				log.Fatal("can't set gauge", err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

func parseUpdate(path string) map[string]string {
	regex := regexp.MustCompile(
		`/update/(?P<type>\w+)(?:/(?P<metrics>\w+))?(?:/(?P<value>[0-9\.\-]+))?`,
	)
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
